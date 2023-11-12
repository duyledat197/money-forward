package http_server

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"user-management/internal/entities"
)

// Middleware represents options that can be used to configure http server
type Middleware interface {
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

// corsMiddleware represents option that allow cors (Cross-origin resource sharing).
type corsMiddleware struct {
	allowMethods []string
}

// Wrap is an implementation of [corsOptions] to wrap next handler into cors handler.
func (m *corsMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.allowMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization")
		if r.Method != "OPTIONS" {
			next.ServeHTTP(w, r)
		}
	})
}

func WithCors(methods ...string) Middleware {
	var allowMethods []string
	if len(methods) == 0 {
		allowMethods = DefaultAllowMethods
	}
	return &corsMiddleware{
		allowMethods: allowMethods,
	}
}

type rbacMiddleware struct {
	rbacMap map[string][]entities.User_Role
}

func (m *rbacMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := joinPath(r.Method, r.URL.Path)
		validRoles, ok := m.rbacMap[route]
		if !ok {
			next.ServeHTTP(w, r)
		}
		info, err := ExtractUserInfoFromContext(r.Context())
		if err != nil {
			errorResponse(w, http.StatusUnauthorized, fmt.Errorf("authorization is not valid"))
			return
		}

		if !slices.Contains(validRoles, entities.User_Role(info.Role)) {
			errorResponse(w, http.StatusForbidden, fmt.Errorf("authorization is not valid"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithRBAC(rbacMap map[string][]entities.User_Role) Middleware {
	return &rbacMiddleware{
		rbacMap: rbacMap,
	}
}
