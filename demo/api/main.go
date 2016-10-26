package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/coreos/dex/pkg/log"
	"github.com/julienschmidt/httprouter"
)

// CONSTANTS
const (
	// Environment Vars
	HOST = "APIHOST"
	PORT = "APIPORT"

	// HTTP Constants
	RESOURCE_ID = "id"
)

type Resource struct {
	Id       string `json:"id, omitempty"`
	Resource string `json:"resource, omitempty"`
}

var resources map[string]string

func HandleAddResource(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := &Resource{}
	err := processHttpRequest(r, request)
	var response *Resource
	if err == nil {
		resources[request.Id] = request.Resource
		response = request
	}
	processHttpResponse(w, response, err, http.StatusOK)
}

func HandleGetResource(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	var response *Resource
	var statusCode int
	if val, ok := resources[ps.ByName(RESOURCE_ID)]; ok {
		response = &Resource{
			Id:       ps.ByName(RESOURCE_ID),
			Resource: val,
		}
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusNotFound
	}
	processHttpResponse(w, response, nil, statusCode)
}

func HandlePutResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	request := &Resource{}
	err := processHttpRequest(r, request)
	var response *Resource
	if err == nil {
		id := ps.ByName("id")
		resources[id] = request.Resource
		response = request
	}
	processHttpResponse(w, response, err, http.StatusOK)
}

func HandleDelResource(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	var statusCode int
	if _, ok := resources[ps.ByName(RESOURCE_ID)]; ok {
		delete(resources, ps.ByName(RESOURCE_ID))
		statusCode = http.StatusNoContent
	} else {
		statusCode = http.StatusNotFound
	}
	processHttpResponse(w, nil, nil, statusCode)
}

func HandleListResources(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	response := make([]Resource, len(resources))
	for key, val := range resources {
		response = append(response, Resource{Id: key, Resource: val})
	}
	processHttpResponse(w, response, nil, http.StatusOK)
}

func main() {

	// Create the muxer to handle the actual endpoints
	router := httprouter.New()

	router.POST("/resources", HandleAddResource)
	router.GET("/resources/:id", HandleGetResource)
	router.PUT("/resources/:id", HandlePutResource)
	router.DELETE("/resources/:id", HandleDelResource)
	router.GET("/resources", HandleListResources)

	// Start server
	resources = make(map[string]string)
	host := os.Getenv(HOST)
	port := os.Getenv(PORT)
	log.Fatal(http.ListenAndServe(host+":"+port, router))

}

// Private Helper Methods

func processHttpRequest(r *http.Request, request interface{}) error {
	// Decode request if passed
	if request != nil {
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return err
		}
	}

	return nil
}

func processHttpResponse(w http.ResponseWriter, response interface{}, err error, responseCode int) {
	if err != nil {
		http.Error(w, err.Error(), responseCode)
		return
	}

	var data []byte
	if response != nil {
		data, err = json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	switch responseCode {
	case http.StatusOK:
		w.Write(data)
	case http.StatusCreated:
		w.Write(data)
	}

	w.WriteHeader(responseCode)

}
