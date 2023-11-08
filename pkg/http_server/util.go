package http_server

import (
	"fmt"
	"strings"
)

const slash = "/"

func joinPath(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}

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
