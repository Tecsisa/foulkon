package main

import (
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/tecsisa/authorizr/authorizr"
	internalhttp "github.com/tecsisa/authorizr/http"
	"net/http"
	"os"
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
		return
	}

	proxy.Logger.Printf("Server running in %v:%v", proxy.Host, proxy.Port)
	if proxy.CertFile != "" && proxy.KeyFile != "" {
		proxy.Logger.Fatal(http.ListenAndServeTLS(proxy.Host+":"+proxy.Port, proxy.CertFile, proxy.KeyFile, internalhttp.ProxyHandlerRouter(proxy)).Error())
	} else {
		proxy.Logger.Fatal(http.ListenAndServe(proxy.Host+":"+proxy.Port, internalhttp.ProxyHandlerRouter(proxy)).Error())
	}
}
