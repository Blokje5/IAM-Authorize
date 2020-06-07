package middleware

import (
	"net/http"

	"github.com/blokje5/iam-server/pkg/server"
)

type errorHandler struct {
	h HandlerWithError
}

func NewErrorHandler(handler HandlerWithError) *errorHandler {
	return &errorHandler{
		h: handler,
	}
}

func (e *errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.h(w, r)

	if res, ok := err.(*server.ErrorResponse); ok == true {
		http.Error(w, res.Error(), res.Status)
	} else {
		http.Error(w, server.NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
	}
}
