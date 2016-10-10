package xrequestid

import (
	"net/http"

	"github.com/Tecsisa/foulkon/middleware"
	"github.com/satori/go.uuid"
)

// XRequestId middleware system
type XRequestIdMiddleware struct{}

// NewXRequestIdMiddleware returns a configured XRequestIdMiddleware
func NewXRequestIdMiddleware() *XRequestIdMiddleware {
	return &XRequestIdMiddleware{}
}

func (r *XRequestIdMiddleware) Action(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewV4().String()
		r.Header.Set(middleware.REQUEST_ID_HEADER, requestID)
		w.Header().Add(middleware.REQUEST_ID_HEADER, requestID)
		next.ServeHTTP(w, r)
	})
}

func (rm *XRequestIdMiddleware) GetInfo(r *http.Request, mc *middleware.MiddlewareContext) {
	mc.XRequestId = r.Header.Get(middleware.REQUEST_ID_HEADER)
}
