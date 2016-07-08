package main

import (
	"flag"
	"fmt"
	"net/http"

	"os"

	"github.com/pelletier/go-toml"
	"github.com/tecsisa/authorizr/authorizr"
	internalhttp "github.com/tecsisa/authorizr/http"
)

func main() {
	// Retrieve config file
	fs := flag.NewFlagSet("authorizr", flag.ExitOnError)
	configFile := fs.String("config-file", "", "Config file for Authorizr")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Access to file
	config, err := toml.LoadFile(*configFile)
	if err != nil {
		fmt.Printf("Cannot read configuration file %v, error: %v", *configFile, err)
		os.Exit(1)
	}

	// Create Worker
	core, err := authorizr.NewWorker(config)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	core.Logger.Printf("Server running in %v:%v", core.Host, core.Port)
	if core.CertFile != "" && core.KeyFile != "" {
		core.Logger.Fatal(http.ListenAndServeTLS(core.Host+":"+core.Port, core.CertFile, core.KeyFile, internalhttp.WorkerHandlerRouter(core)).Error())
	} else {
		core.Logger.Fatal(http.ListenAndServe(core.Host+":"+core.Port, internalhttp.WorkerHandlerRouter(core)).Error())
	}
}
