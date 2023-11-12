package http_server

import (
	"fmt"
	"strings"
)

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
		isMatch = isMatch && (bracketRegex.MatchString(el) || el == sourceEls[i])
	}

	return isMatch
}
