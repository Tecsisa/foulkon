package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Meta struct {
	Address    string
	HttpClient http.Client
}

type Command interface {
	Run(args []string) (string, error)
}

// Helper func for updating request params
func (m *Meta) prepareRequest(method string, url string, postContent map[string]string, queryParams map[string]string) (*http.Request, error) {
	url = m.Address + url
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
	resp, err := m.HttpClient.Do(req)
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
