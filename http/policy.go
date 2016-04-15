package http

import (
	"github.com/tecsisa/authorizr/authorizr"
	"io"
	"net/http"
)

func handleGetPolicy(core *authorizr.Core) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			io.WriteString(w, core.GetPolicyAPI().GetPolicies("/mipath"))
		default:
			authorizr.RespondError(w,http.StatusBadRequest,nil)
		}
	})
}