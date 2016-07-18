package main

import (
	"flag"
	"fmt"
	"net/http"

	"os"

	"os/signal"
	"syscall"

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
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sigrecv := <-sig
		switch sigrecv {
		case syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT:
			core.Logger.Infof("Signal '%v' received, closing worker...", sigrecv.String())
			authorizr.CloseWorker()
		default:
			core.Logger.Warnf("Unknown OS signal received, ignoring...")
		}
	}()

	core.Logger.Infof("Server running in %v:%v", core.Host, core.Port)
	if core.CertFile != "" && core.KeyFile != "" {
		core.Logger.Error(http.ListenAndServeTLS(core.Host+":"+core.Port, core.CertFile, core.KeyFile, internalhttp.WorkerHandlerRouter(core)).Error())
	} else {
		core.Logger.Error(http.ListenAndServe(core.Host+":"+core.Port, internalhttp.WorkerHandlerRouter(core)).Error())
	}

	authorizr.CloseWorker()

}
