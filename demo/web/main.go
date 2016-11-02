package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"strconv"

	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-oidc/oidc"
	"github.com/julienschmidt/httprouter"
	"time"
)

// CONSTANTS
const (
	// Environment Vars
	WEBHOST = "WEBHOST"
	WEBPORT = "WEBPORT"
	APIHOST = "APIHOST"
	APIPORT = "APIPORT"

	// Foulkon
	FOULKONHOST = "FOULKON_WORKER_HOST"
	FOULKONPORT = "FOULKON_WORKER_PORT"

	// OIDC
	OIDCCLIENTID     = "OIDC_CLIENT_ID"
	OIDCCLIENTSECRET = "OIDC_CLIENT_SECRET"
	OIDCIDPDISCOVERY = "OIDC_IDP_DISCOVERY"
)

type Node struct {
	// URLs
	WebBaseUrl     string
	APIBaseUrl     string
	FoulkonBaseUrl string

	// Profile
	UserId string
	Token  string
	Roles  []UserGroups

	// Table Resources
	ResourceTableElements *ResourceTableElements

	// Msg
	Message string

	// Error
	HttpErrorStatusCode int
	ErrorMessage        string
}

type Resource struct {
	Id       string `json:"id, omitempty"`
	Resource string `json:"resource, omitempty"`
}

type UserGroups struct {
	Org      string    `json:"org, omitempty"`
	Name     string    `json:"name, omitempty"`
	CreateAt time.Time `json:"joined, omitempty"`
}

type GetGroupsByUserIdResponse struct {
	Groups []UserGroups `json:"groups, omitempty"`
	Limit  int          `json:"limit, omitempty"`
	Offset int          `json:"offset, omitempty"`
	Total  int          `json:"total, omitempty"`
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
var errorTemplate *template.Template
var client = http.DefaultClient
var logger *logrus.Logger
var node = new(Node)

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

	// Get foulkon url
	foulkonhost := os.Getenv(FOULKONHOST)
	foulkonport := os.Getenv(FOULKONPORT)
	node.FoulkonBaseUrl = "http://" + foulkonhost + ":" + foulkonport + "/api/v1"

	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	// Create oidc client
	client, err := createOidcClient(host, port)
	if err != nil {
		logger.Fatalf("There was an error creating oidc client: %v", err)
	}

	if client == nil {
		logger.Fatal("Nil OIDC client")
	}

	// Create templates
	createTemplates()

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
	router.GET("/login", HandleLoginFunc(client))
	router.GET("/callback", HandleCallbackFunc(client))

	// Start server
	logger.Infof("Server running in %v:%v", host, port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// HANDLERS

func HandlePage(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	if node.Token != "" {
		request, err := http.NewRequest(http.MethodGet, node.FoulkonBaseUrl+"/users/"+node.UserId+"/groups", nil)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		request.SetBasicAuth("admin", "admin")

		response, err := client.Do(request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			err := mainTemplate.Execute(w, node)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		buffer := new(bytes.Buffer)
		if _, err := buffer.ReadFrom(response.Body); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := new(GetGroupsByUserIdResponse)
		if err := json.Unmarshal(buffer.Bytes(), &res); err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		node.Roles = res.Groups

	}
	err := mainTemplate.Execute(w, node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleListResources(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	request, err := http.NewRequest(http.MethodGet, node.APIBaseUrl+"/resources", nil)
	if err != nil {
		logger.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.Header.Set("Authorization", "Bearer "+node.Token)

	response, err := client.Do(request)
	if err != nil {
		logger.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		renderErrorTemplate(w, response.StatusCode, fmt.Sprintf("List method error. Error code: %v", http.StatusText(response.StatusCode)))
		return
	}

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(response.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
		request.Header.Set("Authorization", "Bearer "+node.Token)

		response, err := client.Do(request)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if response.StatusCode != http.StatusOK {
			renderErrorTemplate(w, response.StatusCode, fmt.Sprintf("Add method error. Error code: %v", http.StatusText(response.StatusCode)))
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
		request.Header.Set("Authorization", "Bearer "+node.Token)

		response, err := client.Do(request)
		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if response.StatusCode != http.StatusOK {
			renderErrorTemplate(w, response.StatusCode, fmt.Sprintf("Update method not available. Error code: %v", http.StatusText(response.StatusCode)))
			return
		}

		node.Message = "Resource Updated!"
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
		request.Header.Set("Authorization", "Bearer "+node.Token)
		response, err := client.Do(request)

		if err != nil {
			logger.Info(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if response.StatusCode != http.StatusNoContent {
			renderErrorTemplate(w, response.StatusCode, fmt.Sprintf("Remove method not available. Error code: %v", http.StatusText(response.StatusCode)))
			return
		}

		node.Message = "Resource deleted!"
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

func HandleLoginFunc(c *oidc.Client) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		oac, err := c.OAuthClient()
		if err != nil {
			panic("unable to proceed")
		}

		u, err := url.Parse(oac.AuthCodeURL("", "", ""))
		if err != nil {
			panic("unable to proceed")
		}
		http.Redirect(w, r, u.String(), http.StatusFound)
	}
}

func HandleCallbackFunc(c *oidc.Client) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		code := r.URL.Query().Get("code")
		if code == "" {
			renderErrorTemplate(w, http.StatusBadRequest, fmt.Sprint("code query param must be set"))
			return
		}

		tok, err := c.ExchangeAuthCode(code)
		if err != nil {
			renderErrorTemplate(w, http.StatusBadRequest, fmt.Sprintf("unable to verify auth code with issuer: %v", err))
			return
		}

		claims, err := tok.Claims()
		if err != nil {
			renderErrorTemplate(w, http.StatusBadRequest, fmt.Sprintf("unable to construct claims: %v", err))
			return
		}

		node.Token = tok.Encode()
		node.UserId, _, _ = claims.StringClaim("sub")

		http.Redirect(w, r, node.WebBaseUrl, http.StatusFound)
	}
}

// Aux methods

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

	errorTemplate, err = template.ParseGlob("tmpl/error.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}

	errorTemplate, err = errorTemplate.ParseGlob("tmpl/base/*.html")
	if err != nil {
		log.Fatalf("Template can't be parsed: %v", err)
	}
}

func createOidcClient(host, port string) (*oidc.Client, error) {
	// OIDC client basics
	redirectURL := "http://" + host + ":" + port + "/callback"
	discovery := os.Getenv(OIDCIDPDISCOVERY)

	// OIDC client credentials
	cc := oidc.ClientCredentials{
		ID:     os.Getenv(OIDCCLIENTID),
		Secret: os.Getenv(OIDCCLIENTSECRET),
	}

	logger.Infof("Configured OIDC client with values: redirectURL: %v, IDP discovery: %v", redirectURL, discovery)

	var cfg oidc.ProviderConfig
	var err error
	/*for {*/
	cfg, err = oidc.FetchProviderConfig(http.DefaultClient, discovery)
	if err != nil {
		return nil, err
	}
	/*
		sleep := 3 * time.Second
		logger.Infof("failed fetching provider config, trying again in %v: %v", sleep, err)
		time.Sleep(sleep)
	}*/

	logger.Infof("fetched provider config from %s: %#v", discovery, cfg)

	ccfg := oidc.ClientConfig{
		ProviderConfig: cfg,
		Credentials:    cc,
		RedirectURL:    redirectURL,
	}

	oidcClient, err := oidc.NewClient(ccfg)
	if err != nil {
		return nil, err
	}

	oidcClient.SyncProviderConfig(discovery)

	return oidcClient, nil

}

func renderErrorTemplate(w http.ResponseWriter, code int, msg string) {
	node.HttpErrorStatusCode = code
	node.ErrorMessage = msg
	err := errorTemplate.Execute(w, node)
	if err != nil {
		logger.Info(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
