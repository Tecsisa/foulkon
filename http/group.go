package http

import (
	"io"
	"net/http"

	"github.com/tecsisa/authorizr/authorizr"
)

func handleGroups(core *authorizr.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getGroups(core, w, r)
		default:
			core.RespondError(w, http.StatusBadRequest, nil)
		}
	})
}

func getGroups(core *authorizr.Core, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, core.GetGroupAPI().GetGroups("/mipath"))
}
