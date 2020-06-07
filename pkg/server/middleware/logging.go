package middleware

import (
	"net/http"
	"github.com/blokje5/iam-server/pkg/log"
)

type loggingHandler struct {
	h http.Handler
	logger *log.Logger
}

// NewLoggingHandler returns a logging middleware that adheres the the HTTP handler interface
func NewLoggingHandler(handler http.Handler) http.Handler {
	return &loggingHandler{
		h: handler,
		logger: log.GetLogger(),
	}
}

func (l *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	method := r.Method
	path := r.URL.EscapedPath()

	l.logger.Infof("%v - %v %v", ip, method, path)
	l.h.ServeHTTP(w, r)
}