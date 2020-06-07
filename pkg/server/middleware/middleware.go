package middleware

import "net/http"


type Middleware func (http.Handler) http.Handler

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

