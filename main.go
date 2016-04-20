package main

import (
	"fmt"
	"net/http"

	"os"

	"github.com/tecsisa/authorizr/authorizr"
	internalhttp "github.com/tecsisa/authorizr/http"
)

func main() {

	// Retrieve datasource name
	datasourcename := "/tmp/authorizer/sql.bin"
	// Log dir
	logFileDir := "/tmp/authorizer/authorizer.log"
	// Create log file
	logfile, err := os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	// Create a core
	coreconfig := &authorizr.CoreConfig{
		LogFile:        logfile,
		DatasourceName: datasourcename,
	}
	core, err := authorizr.NewCore(coreconfig)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	core.Logger.Printf("Server running - binding :8000")
	core.Logger.Fatal(http.ListenAndServe(":8000", internalhttp.Handler(core)).Error())
}
