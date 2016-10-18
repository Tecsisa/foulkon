package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const (
	BAD_FLAGS        = 1
	NOT_FOUND        = 2
	DEFAULT_ADDRESS  = "http://127.0.0.1:8080"
	FOULKON_API_PATH = "/api/v1/"
)

type Meta struct {
	address    string
	httpClient http.Client
}

type Command interface {
	Run(args []string) (string, error)
}

func main() {

	help := `Foulkon CLI usage: foulkon [-address=http://1.2.3.4:8080] <command> <action> [<args>]

Available commands:
	user
	group
	policy
	authorize

To get more help, please execute this cli with a <command>

`

	userHelp := `user actions:
	get -id=xxx                   retrieve user xxx
	get-all                       retrieve users
	groups -id=xxx                retrieve user's groups
	create -id=xxx -path=/path    create user 'xxx' with path '/path'
	update -id=xxx -path=/new     update user 'xxx' with path '/new'
	delete -id=xxx                create user 'xxx' with path '/path'

`

	var httpClient http.Client
	var address string
	var command Command

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Printf("%s", help)
		os.Exit(1)
	}

	// get foulkon address
	flag.StringVar(&address, "address", DEFAULT_ADDRESS, "Foulkon Worker address")
	flag.Parse()

	// remove address flag
	r, _ := regexp.Compile(`address`)
	for i, arg := range args {
		if arg == "-address" || arg == "--address" {
			args = append(args[:i], args[i+2:]...)
			break
		}
		if r.MatchString(arg) {
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	var meta Meta = Meta{
		address:    address,
		httpClient: httpClient,
	}

	// force -h flag
	if len(args) < 2 {
		args = append(args, "-h")
	}

	// parse command/action
	switch args[0] {
	case "user":
		switch args[1] {
		case "get":
			command = &GetUserCommand{meta}
		case "get-all":
			command = &GetAllUsersCommand{meta}
		case "groups":
			command = &GetUserGroupsCommand{meta}
		case "create":
			command = &CreateUserCommand{meta}
		case "delete":
			command = &DeleteUserCommand{meta}
		case "-h":
			fmt.Printf(userHelp)
			os.Exit(1)
		default:
			fmt.Printf(userHelp)
			os.Exit(1)
		}
	default:
		fmt.Printf("%s", help)
		os.Exit(1)
	}
	msg, err := command.Run(args[2:])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println(msg)
	os.Exit(0)
}

// Helper func for updating request params
func (m *Meta) prepareRequest(method string, url string, postContent map[string]string, queryParams map[string]string) (*http.Request, error) {
	url = m.address + url
	// insert post content to body
	var body *bytes.Buffer
	if postContent != nil {
		payload, err := json.Marshal(postContent)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(payload)
	}
	if body == nil {
		body = bytes.NewBuffer([]byte{})
	}
	// initialize http request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// add basic auth
	req.SetBasicAuth("admin", "admin")

	// add query params
	if queryParams != nil {
		values := req.URL.Query()
		for k, v := range queryParams {
			values.Add(k, v)
		}
		req.URL.RawQuery = values.Encode()
	}

	return req, nil
}

func (m *Meta) makeRequest(req *http.Request) (string, error) {
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	// read body
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(resp.Body)
	// json pretty-print
	var out bytes.Buffer
	err = json.Indent(&out, buffer.Bytes(), "", "\t")
	if err != nil {
		return "", err
	}
	msg := out.String()
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		if msg == "" {
			msg = "Operation succeeded"
		}
		return msg, nil
	default:
		return "", errors.New("Operation failed, received HTTP status code " + strconv.Itoa(resp.StatusCode))
	}
}
