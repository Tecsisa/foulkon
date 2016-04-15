package main

import (
	"fmt"
	"github.com/tecsisa/authorizr/authorizr"
	internalhttp "github.com/tecsisa/authorizr/http"
	"net/http"
)

func main() {
	core := authorizr.NewCore()
	err := http.ListenAndServe(":8000", internalhttp.Handler(core))
	fmt.Errorf(err.Error())
}
