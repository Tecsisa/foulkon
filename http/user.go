package http

import (
	"net/http"
	"github.com/tecsisa/authorizr/authorizr"
	"encoding/json"
)

func handleGetUsers(core *authorizr.Core) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			b, err := json.Marshal(core.GetUserAPI().GetListUsers("/mipath"))
			if err != nil {
				authorizr.RespondError(w, http.StatusBadRequest, err)
			}

			w.Write(b)
		default:
			authorizr.RespondError(w,http.StatusBadRequest,nil)
		}
	})
}