// Package http_server provides an easy implementation handler
// for http server with register handle by method, path, generic handler
package http_server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"maps"
	"net/http"
	"regexp"
	"strings"

	"user-management/configs"
	"user-management/pkg/reflect_utils"
)

var re = regexp.MustCompile(`\{(.*?)\}`)

type handler[Request, Response any] func(context.Context, *Request) (*Response, error)
type httpHandler func(http.ResponseWriter, *http.Request)

type HttpServer struct {
	*http.ServeMux
	logger     *slog.Logger
	endpoint   *configs.Endpoint
	handlerMap map[string]httpHandler
	server     *http.Server
}

func NewHttpServer(endpoint *configs.Endpoint, logger *slog.Logger) *HttpServer {
	mux := http.NewServeMux()
	return &HttpServer{
		mux,
		logger,
		endpoint,
		make(map[string]httpHandler),
		&http.Server{
			Handler: mux,
			Addr:    endpoint.Address(),
		},
	}
}

// Start will start server and matching with processors pattern
func (s *HttpServer) Start(ctx context.Context) error {
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for route, next := range s.handlerMap {
			_, path, _ := strings.Cut(route, " ")
			if isMatchPath(path, r.URL.Path) {
				next(w, appendWildCardParams(path, r))
				return
			}
		}
		errorResponse(w, http.StatusNotFound, fmt.Errorf("not found"))

	})

	s.logger.InfoContext(ctx, "server listening in", "address", s.endpoint.Address())
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Stop will stop server with graceful shutdown and matching with processors pattern
func (s *HttpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Register will register to http server by method, path and handler with generic handler
func Register[Request, Response any](s *HttpServer, method, path string, handler handler[Request, Response]) {
	switch method {
	case http.MethodOptions:
	case http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut:
		s.handlerMap[joinPath(method, path)] = retrieveRequest(handler)
	default:
		log.Fatalf("unsupported method %s for http server", method)
	}
}

// retrieveRequest returns a handler with marshal all body, query, params from http request to request of generic handler
func retrieveRequest[Request, Response any](handler handler[Request, Response]) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := make(map[string]any)

		// retrieve data from request body with Post, Put methods
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}

		bodyMap := make(map[string]any)
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}
		maps.Copy(params, bodyMap)

		// retrieve data from params (ex: with "/users/{id}" we will got the value of id )
		wildcardParams, ok := ctx.Value(&wildcardParamsKey{}).(map[string]any)
		if !ok {
			errorResponse(w, http.StatusInternalServerError, fmt.Errorf("unable to get wildcard params"))
			return
		}

		maps.Copy(params, wildcardParams)

		// retrieve data from queries params (ex: with /users?name=dat we will got value of name)
		for k, v := range r.URL.Query() {
			switch len(v) {
			case 0:
			case 1:
				params[k] = v[0]
			default:
				params[k] = v
			}
		}

		var req Request
		if err := reflect_utils.ConvertMapToStruct(params, &req); err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		resp, err := handler(ctx, &req)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}

		dataResponse(w, resp)
	}
}
