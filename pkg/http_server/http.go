package http_server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"regexp"

	"user-management/configs"
	"user-management/pkg/reflect_utils"

	"github.com/gorilla/mux"
)

var re = regexp.MustCompile(`\{(.*?)\}`)

type handler[Request, Response any] func(context.Context, *Request) (*Response, error)
type httpHandler func(http.ResponseWriter, *http.Request)

type HttpServer struct {
	*http.ServeMux
	logger     *slog.Logger
	endpoint   *configs.Endpoint
	handlerMap map[string]httpHandler
}

func NewHttpServer(endpoint *configs.Endpoint, logger *slog.Logger) *HttpServer {
	return &HttpServer{
		http.NewServeMux(),
		logger,
		endpoint,
		make(map[string]httpHandler),
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for route, next := range s.handlerMap {
			if isMatchPath(route, joinPath(r.Method, r.URL.Path)) {
				next(w, r)
				return
			}
		}
		errorResponse(w, http.StatusNotFound, fmt.Errorf("not found"))

	})

	s.logger.InfoContext(ctx, "server listening in", "address", s.endpoint.Address())
	if err := http.ListenAndServe(s.endpoint.Address(), s.ServeMux); err != nil {
		return err
	}

	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return nil
}

func Register[Request, Response any](s *HttpServer, method, path string, handler handler[Request, Response]) {
	switch method {
	case http.MethodGet, http.MethodDelete:
		s.handlerMap[joinPath(method, path)] = handleWithoutBody(handler)
	case http.MethodPost, http.MethodPut:
		s.handlerMap[joinPath(method, path)] = handleWithBody(handler)
	}

}

func handleWithoutBody[Request, Response any](handler handler[Request, Response]) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := make(map[string]any)

		wildcardParams := mux.Vars(r)
		for k, v := range wildcardParams {
			params[k] = v
		}

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

func handleWithBody[Request, Response any](handler handler[Request, Response]) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := r.GetBody()
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}
		readReqBody, err := io.ReadAll(body)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}

		params := make(map[string]any)

		wildcardParams := mux.Vars(r)
		for k, v := range wildcardParams {
			params[k] = v
		}

		for k, v := range r.URL.Query() {
			switch len(v) {
			case 0:
			case 1:
				params[k] = v[0]
			default:
				params[k] = v
			}
		}
		bodyMap := make(map[string]any)
		if err := json.Unmarshal(readReqBody, &bodyMap); err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}
		maps.Copy(params, bodyMap)
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

type response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
	Data    any      `json:"data"`
}

func errorResponse(w http.ResponseWriter, code int, err error) {
	resp := &response{
		Code:    code,
		Message: err.Error(),
		Details: []string{},
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func dataResponse(w http.ResponseWriter, data any) {
	resp := &response{
		Data: data,
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
