package authorizr

import (
	"io"
	"os"

	"errors"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/pelletier/go-toml"
)

var proxy_logfile *os.File

// Proxy - Authorize resources using definitions in proxy config file
type Proxy struct {
	// Server config
	Host string
	Port string

	// Worker location
	WorkerHost string

	// TLS configuration
	CertFile string
	KeyFile  string

	// Logger
	Logger *log.Logger

	// API Resources
	APIResources []APIResource
}

// Representation of external API resources to authorize
type APIResource struct {
	Id     string
	Host   string
	Url    string
	Method string
	Urn    string
	Action string
}

func NewProxy(config *toml.TomlTree) (*Proxy, error) {
	// Create logger
	var logOut io.Writer
	var err error
	logOut = os.Stdout
	loggerType := getDefaultValue(config, "logger.type", "Stdout")
	if loggerType == "file" {
		logFileDir := getDefaultValue(config, "logger.file.dir", "/tmp/authorizr.log")
		proxy_logfile, err = os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		logOut = proxy_logfile
	}
	// Loglevel. defaults to INFO
	loglevel, err := log.ParseLevel(getDefaultValue(config, "logger.level", "info"))
	if err != nil {
		loglevel = log.InfoLevel
	}

	logger := &log.Logger{
		Out:       logOut,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     loglevel,
	}
	logger.Infof("Logger type: %v, LogLevel: %v", loggerType, logger.Level.String())

	// API Resources
	resources := []APIResource{}
	// Retrieve resource tree from toml config file
	tree, ok := config.Get("resources").([]*toml.TomlTree)
	if !ok {
		err := errors.New("No resources retrieved from file")
		logger.Error(err)
		return nil, err
	}
	for _, t := range tree {
		resources = append(resources, APIResource{
			Id:     getDefaultValue(t, "id", ""),
			Host:   getDefaultValue(t, "host", ""),
			Url:    getDefaultValue(t, "url", ""),
			Method: getDefaultValue(t, "method", ""),
			Urn:    getDefaultValue(t, "urn", ""),
			Action: getDefaultValue(t, "action", ""),
		})
		logger.Infof("Added resource %v", getDefaultValue(t, "id", ""))
	}

	host, err := getMandatoryValue(config, "server.host")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	port, err := getMandatoryValue(config, "server.port")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	workerHost, err := getMandatoryValue(config, "server.worker-host")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &Proxy{
		Host:         host,
		Port:         port,
		WorkerHost:   workerHost,
		CertFile:     getDefaultValue(config, "server.certfile", ""),
		KeyFile:      getDefaultValue(config, "server.keyfile", ""),
		Logger:       logger,
		APIResources: resources,
	}, nil
}

func CloseProxy() int {
	status := 0
	if err := proxy_logfile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't close logfile: %v", err)
		status = 1
	}
	return status
}
