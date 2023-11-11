package http_server

import (
	"net/http"
	"strings"
)

// Option represents options that can be used to configure http server
type Option interface {
	Wrap(http.Handler) http.Handler
}

var (
	DefaultAllowMethods = []string{
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}
)

// corsOption represents option that allow cors (Cross-origin resource sharing).
type corsOption struct {
	allowMethods []string
}

// Wrap is an implementation of [corsOptions] to wrap next handler into cors handler.
func (o *corsOption) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(o.allowMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization")
		if r.Method != "OPTIONS" {
			next.ServeHTTP(w, r)
		}
	})
}

func WithCors(methods ...string) Option {
	var allowMethods []string
	if len(methods) == 0 {
		allowMethods = DefaultAllowMethods
	}
	return &corsOption{
		allowMethods: allowMethods,
	}
}
