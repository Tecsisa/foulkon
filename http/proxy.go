package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
)

const (
	// Proxy error codes
	INVALID_DEST_HOST_URL = "InvalidDestinationHostURL"
	HOST_UNREACHABLE      = "HostUnreachableError"
	INTERNAL_SERVER_ERROR = "InternalServerError"
	FORBIDDEN_ERROR       = "ForbiddenError"
)

var rUrnParam, _ = regexp.Compile(`\{(\w+)\}`)

func (h *ProxyHandler) HandleRequest(resource authorizr.APIResource) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		requestID := uuid.NewV4().String()
		w.Header().Set(REQUEST_ID_HEADER, requestID)
		// Retrieve parameters to replace in URN
		parameters := getUrnParameters(resource.Urn)
		urn := resource.Urn
		for _, p := range parameters {
			urn = strings.Replace(urn, p[0], ps.ByName(p[1]), -1)
		}
		if workerRequestID, err := h.checkAuthorization(r, urn, resource.Action); err == nil {
			destURL, err := url.Parse(resource.Host)
			if err != nil {
				h.TransactionErrorLog(r, requestID, workerRequestID, fmt.Sprintf("Error creating destination host URL: %v", err.Error()))
				h.RespondInternalServerError(w, getErrorMessage(INVALID_DEST_HOST_URL, "Invalid destination host"))
				return
			}
			r.URL.Host = destURL.Host
			r.URL.Scheme = destURL.Scheme
			// Clean request URI because net/http send method force this
			r.RequestURI = ""
			// Retrieve requested resource
			res, err := h.client.Do(r)
			if err != nil {
				h.TransactionErrorLog(r, requestID, workerRequestID, fmt.Sprintf("Error calling to destination host resource: %v", err.Error()))
				h.RespondInternalServerError(w, getErrorMessage(HOST_UNREACHABLE, "Error calling destination resource"))
				return
			}

			// Copy the response cookies
			for _, cookie := range res.Cookies() {
				http.SetCookie(w, cookie)
			}
			// Copy the response headers from the target server to the proxy response
			for key, values := range res.Header {
				for _, v := range values {
					w.Header().Add(key, v)
				}
			}

			buffer := new(bytes.Buffer)
			if _, err := buffer.ReadFrom(res.Body); err != nil {
				h.TransactionErrorLog(r, requestID, workerRequestID, fmt.Sprintf("Error reading response from destination: %v", err.Error()))
				h.RespondInternalServerError(w, getErrorMessage(INTERNAL_SERVER_ERROR, "Error reading response from destination"))
				return
			}
			w.Write(buffer.Bytes())
			h.TransactionLog(r, requestID, workerRequestID, "Request accepted")
		} else {
			h.TransactionErrorLog(r, requestID, workerRequestID, fmt.Sprintf("Error in authorization: %v", err.Error()))
			apiError := err.(*api.Error)
			switch apiError.Code {
			case FORBIDDEN_ERROR:
				h.RespondForbidden(w, getErrorMessage(FORBIDDEN_ERROR, ""))
			case api.INVALID_PARAMETER_ERROR, api.REGEX_NO_MATCH:
				h.RespondBadRequest(w, getErrorMessage(api.INVALID_PARAMETER_ERROR, "Bad request"))
			default:
				h.RespondInternalServerError(w, getErrorMessage(INTERNAL_SERVER_ERROR, "Internal server error. Contact the administrator"))
			}
			return
		}
	}
}

func (h *ProxyHandler) checkAuthorization(r *http.Request, urn string, action string) (string, error) {
	workerRequestID := "None"
	if !isFullUrn(urn) {
		return workerRequestID,
			getErrorMessage(api.INVALID_PARAMETER_ERROR, fmt.Sprintf("Urn %v is a prefix, it would be a full urn resource", urn))
	}
	if err := api.AreValidResources([]string{urn}); err != nil {
		return workerRequestID, err
	}
	if err := api.AreValidActions([]string{action}); err != nil {
		return workerRequestID, err
	}

	body, err := json.Marshal(AuthorizeResourcesRequest{
		Action:    action,
		Resources: []string{urn},
	})
	if err != nil {
		return workerRequestID, getErrorMessage(api.UNKNOWN_API_ERROR, err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, h.proxy.WorkerHost+RESOURCE_URL, bytes.NewBuffer(body))
	if err != nil {
		return workerRequestID, getErrorMessage(api.UNKNOWN_API_ERROR, err.Error())
	}
	// Add all headers from original request
	req.Header = r.Header
	// Call worker to retrieve authorization
	res, err := h.client.Do(req)
	if err != nil {
		return workerRequestID, getErrorMessage(HOST_UNREACHABLE, err.Error())
	}

	workerRequestID = res.Header.Get(REQUEST_ID_HEADER)

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return workerRequestID, getErrorMessage(FORBIDDEN_ERROR, "Unauthenticated user")
	case http.StatusForbidden:
		return workerRequestID, getErrorMessage(FORBIDDEN_ERROR, fmt.Sprintf("Restricted access to urn %v", urn))
	case http.StatusOK:
		authzResponse := AuthorizeResourcesResponse{}
		err = json.NewDecoder(res.Body).Decode(&authzResponse)
		if err != nil {
			return workerRequestID, getErrorMessage(api.UNKNOWN_API_ERROR, fmt.Sprintf("Error parsing authorizr response %v", err.Error()))
		}

		if len(authzResponse.ResourcesAllowed) < 1 {
			return workerRequestID,
				getErrorMessage(FORBIDDEN_ERROR, fmt.Sprintf("Restricted access to urn %v", urn))
		}

		// Check urns allowed to find target urn
		allowed := false
		for _, allowedRes := range authzResponse.ResourcesAllowed {
			if allowedRes == urn {
				allowed = true
				break
			}
		}

		if !allowed {
			return workerRequestID,
				getErrorMessage(FORBIDDEN_ERROR, fmt.Sprintf("No access for urn %v received from server", urn))
		}

		return workerRequestID, nil
	default:
		return workerRequestID,
			getErrorMessage(INTERNAL_SERVER_ERROR, fmt.Sprintf("There was a problem retrieving authorization, status code %v", res.StatusCode))
	}
}

// Check parameters in URN to replace with URI parameters
func getUrnParameters(urn string) [][]string {
	match := rUrnParam.FindAllStringSubmatch(urn, -1)
	if match != nil && len(match) > 0 {
		return match
	}
	return nil
}

func isFullUrn(resource string) bool {
	return !strings.ContainsAny(resource, "*")
}

func getErrorMessage(errorCode string, message string) *api.Error {
	if message == "" {
		return &api.Error{
			Code:    errorCode,
			Message: "Forbidden resource. If you need access, contact the administrator",
		}
	} else {
		return &api.Error{
			Code:    errorCode,
			Message: message,
		}
	}
}
