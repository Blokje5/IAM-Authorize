package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// New wraps a fn which type is func(w http.ResponseWriter, r *http.Request, next http.Handler) into a Middleware
func New(fn func(http.ResponseWriter, *http.Request, http.Handler)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fn(w, r, next)
		})
	}
}
