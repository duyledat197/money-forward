package http_server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const slash = "/"

// joinPath returns and joining string by method and path with a space between
func joinPath(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}

// isMatchPath returns an bool value with matching from pattern and source by wildcard params.
// ex: "/users/{id}" will match with "/users/123" but "/users" and "/users/123/accounts" not.
func isMatchPath(pattern, source string) bool {
	patternEls := strings.Split(strings.Trim(pattern, slash), slash)
	sourceEls := strings.Split(strings.Trim(source, slash), slash)
	if len(patternEls) != len(sourceEls) {
		return false
	}
	isMatch := true
	for i, el := range patternEls {
		isMatch = isMatch && (re.MatchString(el) || el == sourceEls[i])
	}

	return isMatch
}

type wildcardParamsKey struct{}

// appendWildCardParams will mapping wildcard params in request and append to the request context
func appendWildCardParams(pattern string, r *http.Request) *http.Request {
	result := make(map[string]any)
	patternEls := strings.Split(strings.Trim(pattern, slash), slash)
	sourceEls := strings.Split(strings.Trim(r.URL.Path, slash), slash)
	for i, el := range patternEls {
		if re.MatchString(el) {
			val := strings.TrimLeft(el, "{")
			val = strings.TrimRight(val, "}")
			result[val] = sourceEls[i]
		}
	}

	return r.WithContext(context.WithValue(r.Context(), &wildcardParamsKey{}, result))
}
