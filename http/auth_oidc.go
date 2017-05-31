package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type CreateOidcProviderRequest struct {
	Name        string   `json:"name,omitempty"`
	Path        string   `json:"path,omitempty"`
	IssuerURL   string   `json:"issuerUrl,omitempty"`
	OidcClients []string `json:"clients,omitempty"`
}

type UpdateOidcProviderRequest struct {
	Name        string   `json:"name,omitempty"`
	Path        string   `json:"path,omitempty"`
	IssuerURL   string   `json:"issuerUrl,omitempty"`
	OidcClients []string `json:"clients,omitempty"`
}

// RESPONSES

type ListOidcProvidersResponse struct {
	Providers []string `json:"providers,omitempty"`
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
	Total     int      `json:"total"`
}

// HANDLERS

func (wh *WorkerHandler) HandleAddOidcProvider(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Process request
	request := &CreateOidcProviderRequest{}
	requestInfo, _, apiErr := wh.processHttpRequest(r, w, nil, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call Auth Provider API to create the new OIDC provider
	response, err := wh.worker.AuthOidcAPI.AddOidcProvider(requestInfo, request.Name, request.Path, request.IssuerURL, request.OidcClients)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (wh *WorkerHandler) HandleGetOidcProviderByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call Auth Provider API to get the provider
	response, err := wh.worker.AuthOidcAPI.GetOidcProviderByName(requestInfo, filterData.AuthProviderName)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleListOidcProviders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call Auth Provider API to list the OIDC Providers
	result, total, err := wh.worker.AuthOidcAPI.ListOidcProviders(requestInfo, filterData)
	// Create response
	response := &ListOidcProvidersResponse{
		Providers: result,
		Offset:    filterData.Offset,
		Limit:     filterData.Limit,
		Total:     total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleUpdateOidcProvider(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdateOidcProviderRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call Auth Provider API to update the OIDC Provider
	response, err := wh.worker.AuthOidcAPI.UpdateOidcProvider(requestInfo, filterData.AuthProviderName,
		request.Name, request.Path, request.IssuerURL, request.OidcClients)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleRemoveOidcProvider(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call Auth Provider API to delete the OIDC Provider
	err := wh.worker.AuthOidcAPI.RemoveOidcProvider(requestInfo, filterData.AuthProviderName)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}
