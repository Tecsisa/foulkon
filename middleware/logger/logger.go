package logger

import (
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
)

// Request logger middleware system
type RequestLoggerMiddleware struct{}

// NewRequestLoggerMiddleware returns a configured RequestLoggerMiddleware
func NewRequestLoggerMiddleware() *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{}
}

// Log all request received
func (reqLogger *RequestLoggerMiddleware) Action(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.TransactionRequestLog(r.Header.Get(middleware.REQUEST_ID_HEADER), r.Header.Get(middleware.USER_ID_HEADER), r)
		next.ServeHTTP(w, r)
	})
}

func (reqLogger *RequestLoggerMiddleware) GetInfo(r *http.Request, mc *middleware.MiddlewareContext) {}
