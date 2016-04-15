package main

import (
	"net/http"
	internalhttp "github.com/tecsisa/authorizr/http"
	"github.com/tecsisa/authorizr/authorizr"
	"fmt"
)

func main() {
	core := authorizr.NewCore()
	err := http.ListenAndServe(":8000", internalhttp.Handler(core))
	fmt.Errorf(err.Error())
}
