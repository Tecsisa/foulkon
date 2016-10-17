package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/foulkon"
	internalhttp "github.com/Tecsisa/foulkon/http"
	"github.com/pelletier/go-toml"
)

func main() {
	// Retrieve config file
	fs := flag.NewFlagSet("foulkon", flag.ExitOnError)
	configFile := fs.String("proxy-file", "", "Config file for proxy")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Access to file
	config, err := toml.LoadFile(*configFile)
	if err != nil {
		fmt.Printf("Cannot read proxy file %v, error: %v", *configFile, err)
		os.Exit(1)
	}

	// Create Proxy
	proxy, err := foulkon.NewProxy(config)
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
				api.Log.Infof("Signal '%v' received, closing proxy...", sigrecv.String())
				os.Exit(foulkon.CloseProxy())
			default:
				api.Log.Warnf("Unknown OS signal received, ignoring...")
			}
		}
	}()

	api.Log.Infof("Server running in %v:%v", proxy.Host, proxy.Port)
	ps := internalhttp.NewProxy(proxy)
	ps.Configuration()
	api.Log.Error(ps.Run().Error())

	os.Exit(foulkon.CloseProxy())
}
