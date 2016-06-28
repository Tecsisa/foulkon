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

func (h *ProxyHandler) handleRequest(resource authorizr.APIResource) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		transactionID := uuid.NewV4().String()
		r.Header.Set("Transaction-ID", transactionID)
		w.Header().Set("Transaction-ID", transactionID)
		// Retrieve parameters to replace in URN
		parameters := getUrnParameters(resource.Urn)
		urn := resource.Urn
		for _, p := range parameters {
			urn = strings.Replace(urn, p[0], ps.ByName(p[1]), -1)
		}
		if err := h.checkAuthorization(r, urn, resource.Action); err == nil {
			// Clean request URI because net/http send method force this
			r.RequestURI = ""
			destURL, err := url.Parse(resource.Host)
			if err != nil {
				h.TransactionErrorLog(r, transactionID, fmt.Sprint("Error creating destination host URL: %v", err.Error()))
				http.Error(w, getErrorMessage(transactionID, INVALID_DEST_HOST_URL), http.StatusForbidden)
				return
			}
			r.URL.Host = destURL.Host
			r.URL.Scheme = destURL.Scheme
			// Retrieve requested resource
			res, err := h.client.Do(r)
			if err != nil {
				h.TransactionErrorLog(r, transactionID, fmt.Sprint("Error calling to destination host resource: %v", err.Error()))
				http.Error(w, getErrorMessage(transactionID, DESTINATION_HOST_RESOURCE_CALL_ERROR), http.StatusForbidden)
				return
			}
			res.Write(w)
			h.TransactionLog(r, transactionID, fmt.Sprint("Request accepted"))
		} else {
			h.TransactionErrorLog(r, transactionID, fmt.Sprintf("Error in authorization: %v", err.Error()))
			http.Error(w, getErrorMessage(transactionID, AUTHORIZATION_ERROR), http.StatusForbidden)
			return
		}
	}
}

func (h *ProxyHandler) checkAuthorization(r *http.Request, urn string, action string) error {
	if !isFullUrn(urn) {
		return &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Urn %v is a prefix, it would be a full urn resource", urn),
		}
	}
	if err := api.IsValidResources([]string{urn}); err != nil {
		return err
	}
	if err := api.IsValidAction([]string{action}); err != nil {
		return err
	}

	body, err := json.Marshal(AuthorizeResourcesRequest{
		Action:    action,
		Resources: []string{urn},
	})
	if err != nil {
		return &api.Error{
			Code:    api.UNKNOWN_API_ERROR,
			Message: err.Error(),
		}
	}

	req, err := http.NewRequest(http.MethodPost, h.proxy.WorkerHost+AUTHORIZE_URL, bytes.NewBuffer(body))
	if err != nil {
		return &api.Error{
			Code:    api.UNKNOWN_API_ERROR,
			Message: err.Error(),
		}
	}
	// Add all headers from original request
	req.Header = r.Header
	// Call worker to retrieve authorization
	res, err := h.client.Do(req)
	if err != nil {
		return &api.Error{
			Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
			Message: err.Error(),
		}
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return &api.Error{
			Code:    FORBIDDEN_ERROR,
			Message: fmt.Sprintf("Unauthenticated user"),
		}
	case http.StatusForbidden:
		return &api.Error{
			Code:    FORBIDDEN_ERROR,
			Message: fmt.Sprintf("Restricted access to urn %v", urn),
		}
	case http.StatusOK:
		authzResponse := AuthorizeResourcesResponse{}
		err = json.NewDecoder(res.Body).Decode(&authzResponse)
		if err != nil {
			return &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: err.Error(),
			}
		}

		if len(authzResponse.ResourcesAllowed) < 1 {
			return &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: fmt.Sprintf("Restricted access to urn %v", urn),
			}
		}

		return nil
	default:
		return &api.Error{
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

func getErrorMessage(transactionID string, errorCode string) string {
	return fmt.Sprintf("Forbidden resource. If you need access, contact the administrators."+
		" Transaction id %v. Error code %v", transactionID, errorCode)
}
