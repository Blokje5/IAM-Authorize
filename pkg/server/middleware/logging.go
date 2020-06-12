package middleware

import (
	"net/http"
	"github.com/blokje5/iam-server/pkg/log"
)

// NewLoggingMiddleware returns a logging middleware that can be added to a middleware chain
func NewLoggingMiddleware(logger *log.Logger) Middleware {
	return New(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		ip := r.RemoteAddr
		method := r.Method
		path := r.URL.EscapedPath()
	
		logger.Infof("%v - %v %v", ip, method, path)
		next.ServeHTTP(w, r)
	})
}
