package header

import (
	"fmt"
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/Tecsisa/foulkon/middleware/auth"
)

// HeaderAuthConnector represents a connector that implements interface of auth connector
type HeaderAuthConnector struct {
	header string
}

// InitHeaderConnector initializes Header connector configuration
func InitHeaderConnector(header string) auth.AuthConnector {
	return &HeaderAuthConnector{
		header: header,
	}
}

// Authenticate retrieves auth header from request and returns its value
func (h HeaderAuthConnector) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hdr := r.Header.Get(h.header)
		if hdr == "" {
			apiError := &api.Error{
				Code:    api.AUTHENTICATION_API_ERROR,
				Message: "header authenticator: no auth header found",
			}
			requestID := r.Header.Get(middleware.REQUEST_ID_HEADER)
			api.LogOperationError(requestID, "", apiError)
			http.Error(rw, fmt.Sprintf("Error %v", apiError.Message), http.StatusUnauthorized)
		} else {
			r.Header.Add(middleware.USER_ID_HEADER, hdr)
			next.ServeHTTP(rw, r)
		}
	})
}

// HeaderAuthConnector retrieves user from header
func (h HeaderAuthConnector) RetrieveUserID(r http.Request) string {
	return r.Header.Get(h.header)
}
