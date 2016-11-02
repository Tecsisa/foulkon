package main

import (
	"encoding/json"
	"net/http"
	"os"

	"bytes"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

// CONSTANTS
const (
	// Environment Vars
	HOST        = "APIHOST"
	PORT        = "APIPORT"
	FOULKONHOST = "FOULKON_WORKER_HOST"
	FOULKONPORT = "FOULKON_WORKER_PORT"

	// HTTP Constants
	RESOURCE_ID = "id"
)

type Resource struct {
	Id       string `json:"id, omitempty"`
	Resource string `json:"resource, omitempty"`
}

var resources map[string]string
var logger *logrus.Logger

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

	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	// Startup roles
	createRolesAndPolicies()

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
	logger.Infof("Server running in %v:%v", host, port)
	logger.Fatal(http.ListenAndServe(host+":"+port, router))

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

	w.WriteHeader(responseCode)

	switch responseCode {
	case http.StatusOK:
		w.Write(data)
	case http.StatusCreated:
		w.Write(data)
	}

}

func createRolesAndPolicies() {
	foulkonhost := os.Getenv(FOULKONHOST)
	foulkonport := os.Getenv(FOULKONPORT)
	url := "http://" + foulkonhost + ":" + foulkonport + "/api/v1"

	createGroupFunc := func(name, path string) error {

		type CreateGroupRequest struct {
			Name string `json:"name, omitempty"`
			Path string `json:"path, omitempty"`
		}
		var body *bytes.Buffer
		data := &CreateGroupRequest{
			Name: name,
			Path: path,
		}

		jsonObject, _ := json.Marshal(data)
		body = bytes.NewBuffer(jsonObject)

		req, _ := http.NewRequest(http.MethodPost, url+"/organizations/demo/groups", body)
		req.SetBasicAuth("admin", "admin")

		_, err := http.DefaultClient.Do(req)
		return err
	}

	type Statement struct {
		Effect    string   `json:"effect, omitempty"`
		Actions   []string `json:"actions, omitempty"`
		Resources []string `json:"resources, omitempty"`
	}

	createPolicyAndAttachToGroup := func(name, path, groupName string, statements []Statement) error {
		client := http.DefaultClient
		type CreatePolicyRequest struct {
			Name       string      `json:"name, omitempty"`
			Path       string      `json:"path, omitempty"`
			Statements []Statement `json:"statements, omitempty"`
		}
		var body *bytes.Buffer
		data := &CreatePolicyRequest{
			Name:       name,
			Path:       path,
			Statements: statements,
		}

		jsonObject, _ := json.Marshal(data)
		body = bytes.NewBuffer(jsonObject)

		req, _ := http.NewRequest(http.MethodPost, url+"/organizations/demo/policies", body)
		req.SetBasicAuth("admin", "admin")

		_, err := client.Do(req)
		if err != nil {
			return err
		}

		// Attach
		req, _ = http.NewRequest(http.MethodPost, url+"/organizations/demo/groups/"+groupName+"/policies/"+name, nil)
		req.SetBasicAuth("admin", "admin")

		_, err = client.Do(req)
		if err != nil {
			return err
		}

		return nil
	}

	// Create read role
	err := createGroupFunc("read", "/path/")
	if err != nil {
		logger.Fatal(err)
	}
	statements := []Statement{
		{
			Effect:  "allow",
			Actions: []string{"example:list", "example:get"},
			Resources: []string{
				"urn:ews:foulkon:demo1:resource/list",
				"urn:ews:foulkon:demo1:resource/demoresources/*",
			},
		},
	}

	err = createPolicyAndAttachToGroup("read-policy", "/path/", "read", statements)
	if err != nil {
		logger.Fatal(err)
	}

	// Create read write role
	err = createGroupFunc("read-write", "/path/")
	if err != nil {
		logger.Fatal(err)
	}
	statements2 := []Statement{
		{
			Effect:  "allow",
			Actions: []string{"example:list", "example:get", "example:add", "example:update", "example:delete"},
			Resources: []string{
				"urn:ews:foulkon:demo1:resource/list",
				"urn:ews:foulkon:demo1:resource/demoresources/*",
				"urn:ews:foulkon:demo1:resource/add",
			},
		},
	}

	err = createPolicyAndAttachToGroup("read-write-policy", "/path/", "read-write", statements2)
	if err != nil {
		logger.Fatal(err)
	}

}
