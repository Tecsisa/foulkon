package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"strconv"
)

// CONSTANTS
const (
	// Environment Vars
	WEBHOST = "WEBHOST"
	WEBPORT = "WEBPORT"
	APIHOST = "APIHOST"
	APIPORT = "APIPORT"

	// Operations
	ADD_OPERATION    = "Add"
	GET_OPERATION    = "Get"
	PUT_OPERATION    = "Put"
	DELETE_OPERATION = "Delete"
	LIST_OPERATION   = "List"
)

type Node struct {
	// Operations
	Operation string

	// URLs
	WebBaseUrl string
	APIBaseUrl string

	// Profile
	UserId string
	Token  string
	Roles  []string

	// Table Resources
	ResourceTableElements *ResourceTableElements

	// Msg
	Message string
}

type Resource struct {
	Id       string `json:"id, omitempty"`
	Resource string `json:"resource, omitempty"`
}

type UpdateResource struct {
	Resource string `json:"resource, omitempty"`
}

type ResourceTableElements struct {
	Resources []Resource
}

var mainTemplate *template.Template
var listTemplate *template.Template
var addTemplate *template.Template
var removeTemplate *template.Template
var updateTemplate *template.Template
var client = http.DefaultClient
var logger *logrus.Logger
var node = new(Node)

func createTemplates() {
	var err error
	mainTemplate, err = template.ParseGlob("tmpl/index.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	mainTemplate, err = mainTemplate.ParseGlob("tmpl/base/*.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	listTemplate, err = template.ParseGlob("tmpl/list.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	listTemplate, err = listTemplate.ParseGlob("tmpl/base/*.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	addTemplate, err = template.ParseGlob("tmpl/add.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	addTemplate, err = addTemplate.ParseGlob("tmpl/base/*.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	removeTemplate, err = template.ParseGlob("tmpl/delete.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	removeTemplate, err = removeTemplate.ParseGlob("tmpl/base/*.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	updateTemplate, err = template.ParseGlob("tmpl/update.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	updateTemplate, err = updateTemplate.ParseGlob("tmpl/base/*.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}
}

func HandlePage(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := mainTemplate.Execute(w, node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleListResources(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	logger.Info("Listing resources")
	request, err := http.NewRequest(http.MethodGet, node.APIBaseUrl+"/resources", nil)
	if err != nil {
		logger.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		logger.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(response.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Infof("Values: %v", string(buffer.Bytes()))

	res := make([]Resource, 0)
	if err := json.Unmarshal(buffer.Bytes(), &res); err != nil {
		logger.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	node.ResourceTableElements = &ResourceTableElements{
		Resources: res,
	}

	node.Message = ""
	err = listTemplate.Execute(w, node)
	if err != nil {
		logger.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleAddResource(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == http.MethodPost {
		logger.Info("Adding resource")
		r.ParseForm()
		res := Resource{
			Id:       r.Form.Get("id"),
			Resource: r.Form.Get("resource"),
		}

		jsonObject, err := json.Marshal(&res)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body := bytes.NewBuffer(jsonObject)

		request, err := http.NewRequest(http.MethodPost, node.APIBaseUrl+"/resources", body)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := client.Do(request)

		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if response.StatusCode == http.StatusOK {
			node.Message = "Resource created!"
		} else {
			node.Message = "Error, status code received  " + strconv.Itoa(response.StatusCode)
		}

		addTemplate.Execute(w, node)
	} else {
		node.Message = ""
		err := addTemplate.Execute(w, node)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleUpdateResource(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == http.MethodPost {
		logger.Info("Updating resource")
		r.ParseForm()
		res := UpdateResource{
			Resource: r.Form.Get("resource"),
		}
		id := r.Form.Get("id")

		jsonObject, err := json.Marshal(&res)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body := bytes.NewBuffer(jsonObject)

		request, err := http.NewRequest(http.MethodPut, node.APIBaseUrl+"/resources/"+id, body)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := client.Do(request)

		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if response.StatusCode == http.StatusOK {
			node.Message = "Resource Updated!"
		} else {
			node.Message = "Error, status code received  " + strconv.Itoa(response.StatusCode)
		}

		updateTemplate.Execute(w, node)
	} else {
		node.Message = ""
		err := updateTemplate.Execute(w, node)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleRemoveResource(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == http.MethodPost {
		logger.Info("Removing resource")
		r.ParseForm()
		id := r.Form.Get("id")

		request, err := http.NewRequest(http.MethodDelete, node.APIBaseUrl+"/resources/"+id, nil)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := client.Do(request)

		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if response.StatusCode == http.StatusNoContent {
			node.Message = "Resource deleted!"
		} else {
			node.Message = "Error, status code received  " + strconv.Itoa(response.StatusCode)
		}

		removeTemplate.Execute(w, node)
	} else {
		node.Message = ""
		err := removeTemplate.Execute(w, node)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func main() {
	// Get API Location
	apiHost := os.Getenv(APIHOST)
	apiPort := os.Getenv(APIPORT)
	apiURL := "http://" + apiHost + ":" + apiPort
	node.APIBaseUrl = apiURL

	// Get Web Location
	host := os.Getenv(WEBHOST)
	port := os.Getenv(WEBPORT)
	webURL := "http://" + host + ":" + port
	node.WebBaseUrl = webURL

	// Create templates
	createTemplates()

	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	router := httprouter.New()
	router.GET("/", HandlePage)
	router.POST("/", HandlePage)
	router.GET("/add", HandleAddResource)
	router.POST("/add", HandleAddResource)
	router.GET("/remove", HandleRemoveResource)
	router.POST("/remove", HandleRemoveResource)
	router.GET("/update", HandleUpdateResource)
	router.POST("/update", HandleUpdateResource)
	router.GET("/list", HandleListResources)

	// Start server
	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
