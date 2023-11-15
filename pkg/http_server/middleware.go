package http_server

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"user-management/internal/entities"
	"user-management/pkg/http_server/xcontext"
	"user-management/pkg/logger"
	"user-management/pkg/token_utils"
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
		if r.Method != http.MethodOptions {
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

// corsMiddleware represents option that implements rbac for authorized.
type rbacMiddleware struct {
	rbacMap map[string][]entities.User_Role
}

func (m *rbacMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var validRoles []entities.User_Role

		for route, roles := range m.rbacMap {
			method, path, _ := strings.Cut(route, space)
			// checking path and method is matching with route
			if isMatchPath(path, r.URL.Path) && method == r.Method {
				validRoles = roles
				break
			}
		}

		if len(validRoles) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		info, err := xcontext.ExtractUserInfoFromContext(r.Context())
		if err != nil {
			errorResponse(w, http.StatusUnauthorized, fmt.Errorf("authorization is not valid: user info not valid"))
			return
		}

		if !slices.Contains(validRoles, entities.User_Role(info.Role)) {
			errorResponse(w, http.StatusForbidden, fmt.Errorf("authorization is not valid: role is not valid"))
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

// authenticateMiddleware represents options that implements authenticate for a request.
type authenticateMiddleware struct {
	tokenGenerator token_utils.Authenticator[*xcontext.UserInfo]
	ignoreRoutes   []string
}

func (m *authenticateMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, route := range m.ignoreRoutes {
			method, path, _ := strings.Cut(route, space)
			// checking path and method is matching with route
			if isMatchPath(path, r.URL.Path) && method == r.Method {
				next.ServeHTTP(w, r)
				return
			}
		}

		schema, tkn, ok := strings.Cut(r.Header.Get("Authorization"), space)
		if !ok || strings.ToLower(schema) != "bearer" {
			errorResponse(w, http.StatusForbidden, fmt.Errorf("authorization is not valid: schema must be bearer"))
			return
		}
		payload, err := m.tokenGenerator.Verify(tkn)
		if err != nil {
			errorResponse(w, http.StatusForbidden, err)
			return
		}

		log.Println(*payload)

		next.ServeHTTP(w, r.WithContext(xcontext.ImportUserInfoToContext(r.Context(), payload)))
	})
}

func WithAuthenticate(tokenGenerator token_utils.Authenticator[*xcontext.UserInfo], ignoreRoutes []string) Middleware {
	return &authenticateMiddleware{
		tokenGenerator: tokenGenerator,
		ignoreRoutes:   ignoreRoutes,
	}
}

// recoveryMiddleware represents options that implements recovery a panic occurs in handle flow for a request.
type recoveryMiddleware struct {
	logger logger.Logger
}

func (m *recoveryMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {

				m.logger.Error("http handle was got an error", "err", err) // May be log this error? Send to sentry?
				errorResponse(w, http.StatusInternalServerError, fmt.Errorf("there was an internal server error"))
			}

		}()

		next.ServeHTTP(w, r)
	})
}

func WithRecovery(logger logger.Logger) Middleware {
	return &recoveryMiddleware{
		logger: logger,
	}
}
