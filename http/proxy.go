package http

import (
	"net/http"
	"net/url"

	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
)

const (

	// Proxy error codes
	INVALID_DEST_HOST_URL                = "InvalidDestinationHostURL"
	DESTINATION_HOST_RESOURCE_CALL_ERROR = "DestinationHostResourceCallError"
	AUTHORIZATION_ERROR                  = "AuthorizationError"
	FORBIDDEN_ERROR                      = "ForbiddenError"
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
				h.RespondForbidden(w, getErrorMessage(INVALID_DEST_HOST_URL))
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
				h.RespondForbidden(w, getErrorMessage(DESTINATION_HOST_RESOURCE_CALL_ERROR))
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
			buffer.ReadFrom(res.Body)
			w.Write(buffer.Bytes())
			h.TransactionLog(r, requestID, workerRequestID, "Request accepted")
		} else {
			h.TransactionErrorLog(r, requestID, workerRequestID, fmt.Sprintf("Error in authorization: %v", err.Error()))
			h.RespondForbidden(w, getErrorMessage(AUTHORIZATION_ERROR))
			return
		}
	}
}

func (h *ProxyHandler) checkAuthorization(r *http.Request, urn string, action string) (string, error) {
	workerRequestID := "None"
	if !isFullUrn(urn) {
		return workerRequestID, &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Urn %v is a prefix, it would be a full urn resource", urn),
		}
	}
	if err := api.IsValidResources([]string{urn}); err != nil {
		return workerRequestID, err
	}
	if err := api.IsValidAction([]string{action}); err != nil {
		return workerRequestID, err
	}

	body, err := json.Marshal(AuthorizeResourcesRequest{
		Action:    action,
		Resources: []string{urn},
	})
	if err != nil {
		return workerRequestID, &api.Error{
			Code:    api.UNKNOWN_API_ERROR,
			Message: err.Error(),
		}
	}

	req, err := http.NewRequest(http.MethodPost, h.proxy.WorkerHost+AUTHORIZE_URL, bytes.NewBuffer(body))
	if err != nil {
		return workerRequestID, &api.Error{
			Code:    api.UNKNOWN_API_ERROR,
			Message: err.Error(),
		}
	}
	// Add all headers from original request
	req.Header = r.Header
	// Call worker to retrieve authorization
	res, err := h.client.Do(req)
	if err != nil {
		return workerRequestID, &api.Error{
			Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
			Message: err.Error(),
		}
	}

	workerRequestID = res.Header.Get(REQUEST_ID_HEADER)

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return workerRequestID, &api.Error{
			Code:    FORBIDDEN_ERROR,
			Message: "Unauthenticated user",
		}
	case http.StatusForbidden:
		return workerRequestID, &api.Error{
			Code:    FORBIDDEN_ERROR,
			Message: fmt.Sprintf("Restricted access to urn %v", urn),
		}
	case http.StatusOK:
		authzResponse := AuthorizeResourcesResponse{}
		err = json.NewDecoder(res.Body).Decode(&authzResponse)
		if err != nil {
			return workerRequestID, &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: fmt.Sprintf("Error parsing authorizr response %v", err.Error()),
			}
		}

		if len(authzResponse.ResourcesAllowed) < 1 {
			return workerRequestID, &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: fmt.Sprintf("Restricted access to urn %v", urn),
			}
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
			return workerRequestID, &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: fmt.Sprintf("No access for urn %v received from server", urn),
			}
		}

		return workerRequestID, nil
	default:
		return workerRequestID, &api.Error{
			Code:    AUTHORIZATION_ERROR,
			Message: fmt.Sprintf("There was a problem calling authorization, status code %v", res.StatusCode),
		}
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
	if strings.ContainsAny(resource, "*") {
		return false
	} else {
		return true
	}
}

func getErrorMessage(errorCode string) *api.Error {
	return &api.Error{
		Code:    errorCode,
		Message: "Forbidden resource. If you need access, contact the administrators.",
	}
}
