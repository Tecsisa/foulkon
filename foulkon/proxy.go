package foulkon

import (
	"io"
	"os"

	"errors"

	"fmt"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/Tecsisa/foulkon/api"
	"github.com/pelletier/go-toml"

	"github.com/Tecsisa/foulkon/database/postgresql"
)

var proxyLogfile *os.File

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

	// API
	ProxyApi api.ProxyResourcesAPI

	// Refresh time
	RefreshTime time.Duration
}

func NewProxy(config *toml.TomlTree) (*Proxy, error) {
	// Create logger
	var logOut io.Writer
	var err error
	logOut = os.Stdout
	loggerType := getDefaultValue(config, "logger.type", "Stdout")
	if loggerType == "file" {
		logFileDir := getDefaultValue(config, "logger.file.dir", "/tmp/foulkon.log")
		proxyLogfile, err = os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		logOut = proxyLogfile
	}
	// Loglevel. defaults to INFO
	loglevel, err := logrus.ParseLevel(getDefaultValue(config, "logger.level", "info"))
	if err != nil {
		loglevel = logrus.InfoLevel
	}

	api.Log = &logrus.Logger{
		Out:       logOut,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     loglevel,
	}
	api.Log.Infof("Logger type: %v, LogLevel: %v", loggerType, api.Log.Level.String())

	// Start DB with API
	var prApi api.ProxyAPI

	dbType, err := getMandatoryValue(config, "database.type")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	switch dbType {
	case "postgres": // PostgreSQL DB
		api.Log.Info("Connecting to postgres database")
		dbdsn, err := getMandatoryValue(config, "database.postgres.datasourcename")
		if err != nil {
			api.Log.Error(err)
			return nil, err
		}
		gormDB, err := postgresql.InitDb(dbdsn,
			getDefaultValue(config, "database.postgres.idleconns", "5"),
			getDefaultValue(config, "database.postgres.maxopenconns", "20"),
			getDefaultValue(config, "database.postgres.connttl", "300"),
		)
		if err != nil {
			api.Log.Error(err)
			return nil, err
		}
		db = gormDB.DB()
		api.Log.Info("Connected to postgres database")

		// Create repository
		repoDB := postgresql.PostgresRepo{
			Dbmap: gormDB,
		}
		prApi = api.ProxyAPI{
			ProxyRepo: repoDB,
		}

	default:
		err := errors.New("Unexpected db_type value in configuration file (Maybe it is empty)")
		api.Log.Error(err)
		return nil, err
	}

	host, err := getMandatoryValue(config, "server.host")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	port, err := getMandatoryValue(config, "server.port")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	workerHost, err := getMandatoryValue(config, "server.worker-host")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}

	refresh, err := time.ParseDuration(getDefaultValue(config, "resources.refresh", "10s"))
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}

	return &Proxy{
		Host:        host,
		Port:        port,
		WorkerHost:  workerHost,
		CertFile:    getDefaultValue(config, "server.certfile", ""),
		KeyFile:     getDefaultValue(config, "server.keyfile", ""),
		ProxyApi:    prApi,
		RefreshTime: refresh,
	}, nil
}

func CloseProxy() int {
	status := 0
	if err := db.Close(); err != nil {
		api.Log.Errorf("Couldn't close DB connection: %v", err)
		status = 1
	}
	if proxyLogfile != nil {
		if err := proxyLogfile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't close logfile: %v", err)
			status = 1
		}
	}
	return status
}
