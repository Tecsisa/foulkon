package authorizr

import (
	"io"
	"os"

	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/pelletier/go-toml"
)

// Proxy authorize resources according to definitions in proxy file
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

// Representation of resources with its correspondence resource and action for authorizr server
type APIResource struct {
	// API definition
	Id     string
	Host   string
	Url    string
	Method string
	// Authorization relationship
	Urn    string
	Action string
}

// Create a Proxy using configuration values
func NewProxy(config *toml.TomlTree) (*Proxy, error) {
	// Create logger
	var logOut io.Writer
	logOut = os.Stdout
	loggerType := getDefaultValue(config, "logger.type", "Stdout")
	if loggerType == "file" {
		logFileDir := getDefaultValue(config, "logger.file.dir", "/tmp/authorizr.log")
		file, err := os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		logOut = file
	}
	// Logger level
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
	tree, ok := config.Get("resources").([]*toml.TomlTree)
	if !ok {
		return nil, errors.New("No resources retrieved from file")
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

	// Return created proxy
	return &Proxy{
		Host:         getMandatoryValue(config, "server.host"),
		Port:         getMandatoryValue(config, "server.port"),
		WorkerHost:   getMandatoryValue(config, "server.worker-host"),
		CertFile:     getDefaultValue(config, "server.certfile", ""),
		KeyFile:      getDefaultValue(config, "server.keyfile", ""),
		Logger:       logger,
		APIResources: resources,
	}, nil
}
