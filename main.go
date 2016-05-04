package main

import (
	"flag"
	"fmt"
	"net/http"

	"os"

	"encoding/json"
	"github.com/tecsisa/authorizr/authorizr"
	internalhttp "github.com/tecsisa/authorizr/http"
	"io/ioutil"
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
	file, e := ioutil.ReadFile(*configFile)
	if e != nil {
		fmt.Printf("Cannot read configuration file %v File error: %v\n", *configFile, e)
		os.Exit(1)
	}

	// Transform json
	var coreConfig *authorizr.CoreConfig
	err := json.Unmarshal(file, &coreConfig)
	if err != nil {
		fmt.Printf("Unmarshal file error: %v\n", err)
		os.Exit(1)
	}

	// Create core
	core, err := authorizr.NewCore(coreConfig)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	core.Logger.Printf("Server running - binding :8000")
	core.Logger.Fatal(http.ListenAndServe(":8000", internalhttp.Handler(core)).Error())
}
