package main

import (
	"flag"
	"fmt"

	"os"

	"os/signal"
	"syscall"

	"github.com/Tecsisa/foulkon/foulkon"
	internalhttp "github.com/Tecsisa/foulkon/http"
	"github.com/pelletier/go-toml"
)

func main() {
	// Retrieve config file
	fs := flag.NewFlagSet("foulkon", flag.ExitOnError)
	configFile := fs.String("config-file", "", "Config file for worker")

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
	core, err := foulkon.NewWorker(config)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for {
			sigrecv := <-sig
			switch sigrecv {
			case syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT:
				core.Logger.Infof("Signal '%v' received, closing worker...", sigrecv.String())
				os.Exit(foulkon.CloseWorker())
			default:
				core.Logger.Warnf("Unknown OS signal received, ignoring...")
			}
		}
	}()

	core.Logger.Infof("Server running in %v:%v", core.Host, core.Port)
	ws := internalhttp.NewWorker(core, internalhttp.WorkerHandlerRouter(core))
	ws.Configuration()
	core.Logger.Error(ws.Run().Error())

	os.Exit(foulkon.CloseWorker())
}
