package http

import (
	"io"
	"net/http"

	"github.com/tecsisa/authorizr/authorizr"
)

func handlePolicies(core *authorizr.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			io.WriteString(w, core.Policyapi.GetPolicies("/mipath"))
		default:
			core.RespondError(w, http.StatusBadRequest, nil)
		}
	})
}
