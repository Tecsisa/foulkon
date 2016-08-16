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
	proxy, err := authorizr.NewProxy(config)
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
			proxy.Logger.Infof("Signal '%v' received, closing proxy...", sigrecv.String())
			authorizr.CloseProxy()
		default:
			proxy.Logger.Warnf("Unknown OS signal received, ignoring...")
		}
	}()

	proxy.Logger.Infof("Server running in %v:%v", proxy.Host, proxy.Port)
	if proxy.CertFile != "" && proxy.KeyFile != "" {
		proxy.Logger.Error(http.ListenAndServeTLS(proxy.Host+":"+proxy.Port, proxy.CertFile, proxy.KeyFile, internalhttp.ProxyHandlerRouter(proxy)).Error())
	} else {
		proxy.Logger.Error(http.ListenAndServe(proxy.Host+":"+proxy.Port, internalhttp.ProxyHandlerRouter(proxy)).Error())
	}

	os.Exit(authorizr.CloseProxy())

}
