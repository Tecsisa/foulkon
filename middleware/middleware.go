package middleware

import "net/http"

const (
	// HTTP Header
	REQUEST_ID_HEADER = "X-Request-Id"
	USER_ID_HEADER    = "X-FOULKON-USER-ID"

	// Middleware names
	AUTHENTICATOR_MIDDLEWARE  = "AUTHENTICATOR"
	XREQUESTID_MIDDLEWARE     = "XREQUESTID"
	REQUEST_LOGGER_MIDDLEWARE = "REQUEST-LOGGER"
)

// MiddlewareHandler handles the HTTP request and applies its list of middlewares before calling the API
type MiddlewareHandler struct {
	Middlewares map[string]Middleware
}

// MiddlewareContext struct contains all parameters used in the context of middlewares
type MiddlewareContext struct {
	// Authenticator middleware
	UserId string
	Admin  bool

	// X-Request-Id middleware
	XRequestId string
}

// Middleware interface with operations that all middlewares must implement
type Middleware interface {
	// Action to apply to each request
	Action(next http.Handler) http.Handler

	// Additional info that middleware use
	GetInfo(r *http.Request, mc *MiddlewareContext)
}

// Handle method execute middlewares in correct order before API handler
func (mwh *MiddlewareHandler) Handle(apiHandler http.Handler) http.Handler {
	var handler http.Handler
	// Wrap target handler to use middleware
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiHandler.ServeHTTP(w, r)
	})
	// Middleware execution order is upside-down because when adding an action
	// it executes itself first, then the rest of the old handler.
	if val, ok := mwh.Middlewares[REQUEST_LOGGER_MIDDLEWARE]; ok {
		handler = val.Action(handler)
	}
	if val, ok := mwh.Middlewares[AUTHENTICATOR_MIDDLEWARE]; ok {
		handler = val.Action(handler)
	}
	if val, ok := mwh.Middlewares[XREQUESTID_MIDDLEWARE]; ok {
		handler = val.Action(handler)
	}

	return handler
}

// GetMiddlewareContext method retrieves all information about middleware context applied to request
func (mwh *MiddlewareHandler) GetMiddlewareContext(r *http.Request) *MiddlewareContext {
	context := new(MiddlewareContext)
	for _, m := range mwh.Middlewares {
		m.GetInfo(r, context)
	}

	return context
}
