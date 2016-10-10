package logger

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/Tecsisa/foulkon/middleware"
)

// Request logger middleware system
type RequestLoggerMiddleware struct {
	log *logrus.Logger
}

// NewRequestLoggerMiddleware returns a configured RequestLoggerMiddleware
func NewRequestLoggerMiddleware(log *logrus.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		log: log,
	}
}

// Log all request received
func (reqLogger *RequestLoggerMiddleware) Action(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: X-Forwarded headers?
		//for header, _ := range r.Header {
		//	println(header, ": ", r.Header.Get(header))
		//}

		reqLogger.log.WithFields(logrus.Fields{
			"requestID": r.Header.Get(middleware.REQUEST_ID_HEADER),
			"method":    r.Method,
			"URI":       r.RequestURI,
			"address":   r.RemoteAddr,
			"user":      r.Header.Get(middleware.USER_ID_HEADER),
		}).Info("")
		next.ServeHTTP(w, r)
	})
}

func (reqLogger *RequestLoggerMiddleware) GetInfo(r *http.Request, mc *middleware.MiddlewareContext) {}
