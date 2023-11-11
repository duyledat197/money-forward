package http_server

import "net/http"

type Option interface {
	Wrap(http.Handler) http.Handler
}

type corsOption struct {
}

func (o *corsOption) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)
	})
}

func WithCors() Option {
	return &corsOption{}
}
