package main

import (
	"fmt"
	"github.com/tecsisa/authorizr/authorizr"
	internalhttp "github.com/tecsisa/authorizr/http"
	"net/http"
)

func main() {
	core := authorizr.NewCore()
	fmt.Printf("Server running - binding :8000")
	err := http.ListenAndServe(":8000", internalhttp.Handler(core))
	fmt.Errorf(err.Error())
}
