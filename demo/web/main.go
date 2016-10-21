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
}

type Resource struct {
	Id       string `json:"id, omitempty"`
	Resource string `json:"resource, omitempty"`
}

type ResourceTableElements struct {
	Resources []Resource
}

var mainTemplate *template.Template
var client = http.DefaultClient
var logger *logrus.Logger
var node = new(Node)

func HandlePage(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := mainTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ListResources(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	logger.Info("LISTING RESOURCES")
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

	node.ResourceTableElements.Resources = res

	err = mainTemplate.Execute(w, node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Create template
	var err error
	mainTemplate, err = template.ParseGlob("tmpl/index.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	mainTemplate, err = mainTemplate.ParseGlob("tmpl/base/*.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	router := httprouter.New()
	router.GET("/", HandlePage)
	router.POST("/", HandlePage)

	// Start server
	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
