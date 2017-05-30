package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

const (
	// Proxy error codes
	INVALID_DEST_HOST_URL = "InvalidDestinationHostURL"
	HOST_UNREACHABLE      = "HostUnreachableError"
	INTERNAL_SERVER_ERROR = "InternalServerError"
	BAD_REQUEST           = "BadRequest"
	FORBIDDEN_ERROR       = "ForbiddenError"
)

// REQUESTS

type CreateProxyResourceRequest struct {
	Name     string             `json:"name,omitempty"`
	Path     string             `json:"path,omitempty"`
	Resource api.ResourceEntity `json:"resource,omitempty"`
}

type UpdateProxyResourceRequest struct {
	Name     string             `json:"name,omitempty"`
	Path     string             `json:"path,omitempty"`
	Resource api.ResourceEntity `json:"resource,omitempty"`
}

// RESPONSES

type ProxyResources struct {
	Resources []api.ProxyResource `json:"resources,omitempty"`
}

type ListProxyResourcesResponse struct {
	Resources []string `json:"resources,omitempty"`
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
	Total     int      `json:"total"`
}

var rUrnParam, _ = regexp.Compile(`\{(\w+)\}`)

func (ph *ProxyHandler) HandleRequest(proxyResource api.ProxyResource) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		requestID := uuid.NewV4().String()
		w.Header().Set(middleware.REQUEST_ID_HEADER, requestID)
		// Retrieve parameters to replace in URN
		parameters := getUrnParameters(proxyResource.Resource.Urn)
		urn := proxyResource.Resource.Urn
		for _, p := range parameters {
			urn = strings.Replace(urn, p[0], ps.ByName(p[1]), -1)
		}
		if workerRequestID, err := ph.checkAuthorization(r, urn, proxyResource.Resource.Action); err == nil {
			destURL, err := url.Parse(proxyResource.Resource.Host)
			if err != nil {
				apiErr := getErrorMessage(INVALID_DEST_HOST_URL, fmt.Sprintf("Error creating destination host URL: %v", err.Error()))
				api.TransactionProxyErrorLogWithStatus(requestID, workerRequestID, r, http.StatusInternalServerError, apiErr)
				WriteHttpResponse(r, w, requestID, "", http.StatusInternalServerError, getErrorMessage(INVALID_DEST_HOST_URL, "Error creating destination host"))
				return
			}
			r.URL.Host = destURL.Host
			r.URL.Scheme = destURL.Scheme
			// Clean request URI because net/http send method force this
			r.RequestURI = ""
			// Retrieve requested resource
			res, err := ph.client.Do(r)
			if err != nil {
				apiErr := getErrorMessage(HOST_UNREACHABLE, fmt.Sprintf("Error calling to destination host resource: %v", err.Error()))
				api.TransactionProxyErrorLogWithStatus(requestID, workerRequestID, r, http.StatusInternalServerError, apiErr)
				WriteHttpResponse(r, w, requestID, "", http.StatusInternalServerError, getErrorMessage(HOST_UNREACHABLE, "Error calling destination resource"))
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

			defer res.Body.Close()
			buffer := new(bytes.Buffer)
			if _, err := buffer.ReadFrom(res.Body); err != nil {
				apiErr := getErrorMessage(INTERNAL_SERVER_ERROR, fmt.Sprintf("Error reading response from destination: %v", err.Error()))
				api.TransactionProxyErrorLogWithStatus(requestID, workerRequestID, r, http.StatusInternalServerError, apiErr)
				WriteHttpResponse(r, w, requestID, "", http.StatusInternalServerError, apiErr)
				return
			}
			w.WriteHeader(res.StatusCode)
			w.Write(buffer.Bytes())
			api.TransactionProxyLog(requestID, workerRequestID, r, "Request accepted")
		} else {
			apiError := err.(*api.Error)
			var statusCode int
			var responseErr *api.Error
			switch apiError.Code {
			case FORBIDDEN_ERROR:
				statusCode = http.StatusForbidden
				responseErr = getErrorMessage(FORBIDDEN_ERROR, "")
			case api.INVALID_PARAMETER_ERROR, api.REGEX_NO_MATCH, BAD_REQUEST:
				statusCode = http.StatusBadRequest
				responseErr = getErrorMessage(api.INVALID_PARAMETER_ERROR, "Bad request")
			default:
				statusCode = http.StatusInternalServerError
				responseErr = getErrorMessage(INTERNAL_SERVER_ERROR, "Internal server error. Contact the administrator")
			}
			WriteHttpResponse(r, w, requestID, "", statusCode, responseErr)
			api.TransactionProxyErrorLogWithStatus(requestID, workerRequestID, r, statusCode, apiError)
			return
		}
	}
}

// HANDLERS

func (wh *WorkerHandler) HandleAddProxyResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &CreateProxyResourceRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call proxy Resource API to create proxyResource
	response, err := wh.worker.ProxyApi.AddProxyResource(requestInfo, request.Name, filterData.Org, request.Path, request.Resource)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (wh *WorkerHandler) HandleGetProxyResourceByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call policy API to retrieve policy
	response, err := wh.worker.ProxyApi.GetProxyResourceByName(requestInfo, filterData.Org, filterData.ProxyResourceName)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleListProxyResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call proxy Resource API to create proxyResource
	result, total, err := wh.worker.ProxyApi.ListProxyResources(requestInfo, filterData)
	proxyResources := []string{}
	for _, proxyResource := range result {
		proxyResources = append(proxyResources, proxyResource.Name)
	}
	// Create response
	response := &ListProxyResourcesResponse{
		Resources: proxyResources,
		Offset:    filterData.Offset,
		Limit:     filterData.Limit,
		Total:     total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleUpdateProxyResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdateProxyResourceRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call proxy resource API to update proxy resource
	response, err := wh.worker.ProxyApi.UpdateProxyResource(requestInfo, filterData.Org, filterData.ProxyResourceName, request.Name, request.Path, request.Resource)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleRemoveProxyResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call proxy resource API to remove proxy resource
	err := wh.worker.ProxyApi.RemoveProxyResource(requestInfo, filterData.Org, filterData.ProxyResourceName)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (ph *ProxyHandler) checkAuthorization(r *http.Request, urn string, action string) (string, error) {
	workerRequestID := "None"
	if !isFullUrn(urn) {
		return workerRequestID,
			getErrorMessage(api.INVALID_PARAMETER_ERROR, fmt.Sprintf("Urn %v is a prefix, it would be a full urn resource", urn))
	}
	if err := api.AreValidResources([]string{urn}, api.RESOURCE_EXTERNAL); err != nil {
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

	req, err := http.NewRequest(http.MethodPost, ph.proxy.WorkerHost+RESOURCE_URL, bytes.NewBuffer(body))
	if err != nil {
		return workerRequestID, getErrorMessage(api.UNKNOWN_API_ERROR, err.Error())
	}
	// Add all headers from original request
	req.Header = r.Header
	// Call worker to retrieve authorization
	res, err := ph.client.Do(req)
	if err != nil {
		return workerRequestID, getErrorMessage(HOST_UNREACHABLE, err.Error())
	}

	defer res.Body.Close()

	workerRequestID = res.Header.Get(middleware.REQUEST_ID_HEADER)

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return workerRequestID, getErrorMessage(FORBIDDEN_ERROR, "Unauthenticated user")
	case http.StatusForbidden:
		return workerRequestID, getErrorMessage(FORBIDDEN_ERROR, fmt.Sprintf("Restricted access to urn %v", urn))
	case http.StatusBadRequest:
		return workerRequestID, getErrorMessage(BAD_REQUEST, "Invalid request")
	case http.StatusOK:
		authzResponse := AuthorizeResourcesResponse{}
		err = json.NewDecoder(res.Body).Decode(&authzResponse)
		if err != nil {
			return workerRequestID, getErrorMessage(api.UNKNOWN_API_ERROR, fmt.Sprintf("Error parsing foulkon response %v", err.Error()))
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
	}
	return &api.Error{
		Code:    errorCode,
		Message: message,
	}
}
