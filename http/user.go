package http

import (
	"net/http"
	"github.com/tecsisa/authorizr/authorizr"
	"io"
)

func handleGetUsers(core *authorizr.Core) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			io.WriteString(w, core.GetUserAPI().GetListUsers("/mipath"))
		default:
			authorizr.RespondError(w,http.StatusBadRequest,nil)
		}
	})
}