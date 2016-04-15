package http

import (
	"github.com/tecsisa/authorizr/authorizr"
	"io"
	"net/http"
)

func handleGetGroups(core *authorizr.Core) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getGroups(core, w, r)
		default:
			authorizr.RespondError(w,http.StatusBadRequest,nil)
		}
	})
}

func getGroups(core *authorizr.Core, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, core.GetGroupAPI().GetGroups("/mipath"))
}