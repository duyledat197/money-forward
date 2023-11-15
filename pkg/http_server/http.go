// Package http_server provides some quick implementation handler used http handler.
package http_server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"slices"
	"strings"

	"user-management/configs"
	"user-management/pkg/logger"
	"user-management/pkg/reflect_utils"
)

// handler is a presentation for a implementation of a delivery API.
// handler returns a response or error with passing context and request in parameters.
type handler[Request, Response any] func(context.Context, *Request) (*Response, error)

// httpHandler is a presentation for handle func of [net/http]
type httpHandler func(http.ResponseWriter, *http.Request)

// HttpServer represents a http server include [net/http.ServeMux], [user-management/Logger]
type HttpServer struct {
	logger      logger.Logger
	endpoint    *configs.Endpoint
	handlerMap  map[string]httpHandler
	server      *http.Server
	middlewares []Middleware
}

// NewHttpServer returns a custom http server, used to serve a http server.
func NewHttpServer(
	endpoint *configs.Endpoint,
	logger logger.Logger,
	middlewares ...Middleware,
) *HttpServer {
	return &HttpServer{
		logger:      logger,
		endpoint:    endpoint,
		handlerMap:  make(map[string]httpHandler),
		middlewares: middlewares,
	}
}

// Start will start server and matching with processors pattern
func (s *HttpServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc(slash, func(w http.ResponseWriter, r *http.Request) {
		for route, next := range s.handlerMap {
			method, path, _ := strings.Cut(route, space)
			// checking path and method is matching with route
			if isMatchPath(path, r.URL.Path) && method == r.Method {
				next(w, appendWildCardParams(path, r))
				return
			}
		}
		errorResponse(w, http.StatusNotFound, fmt.Errorf("not found"))
	})

	var handler http.Handler = mux
	slices.Reverse(s.middlewares)

	// Merge all middleware handlers into one that can using for register to http server.
	for _, middleware := range s.middlewares {
		handler = middleware.Wrap(handler)
	}

	s.logger.Info("server listening in", "address", s.endpoint.Address())
	if err := http.ListenAndServe(s.endpoint.Address(), handler); err != nil {
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
	case
		http.MethodGet,
		http.MethodDelete,
		http.MethodPost,
		http.MethodPut:
		s.handlerMap[joinPath(method, path)] = handleRequest(handler)
	default:
		log.Fatalf("unsupported method %s for http server", method)
	}
}

// handleRequest returns a handler with marshal all body, query, params
// from http request to request of generic handler.
func handleRequest[Request, Response any](handler handler[Request, Response]) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params, err := retrieveDataFromRequest(w, r)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}

		var req Request
		// convert all params into request struct
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

// retrieveDataFromRequest returns a map that is all query params and body converted from request.
func retrieveDataFromRequest(w http.ResponseWriter, r *http.Request) (map[string]any, error) {
	ctx := r.Context()
	params := make(map[string]any)

	// retrieve data from request body with Post, Put methods
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		bodyMap := make(map[string]any)
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			return nil, err
		}
		maps.Copy(params, bodyMap)
	}

	// retrieve data from wildcard params (ex: with "/users/{id}" we will got the value of id )
	wildcardParams, ok := ctx.Value(&wildcardParamsKey{}).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unable to get wildcard params")
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

	return params, nil
}
