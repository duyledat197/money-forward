package http_server

import (
	"path/filepath"
	"regexp"
)

const (
	slash        = string(filepath.Separator)
	space        = " "
	openBracket  = "{"
	closeBracket = "}"
)

var bracketRegex = regexp.MustCompile(`\{(.*?)\}`)

type (
	wildcardParamsKey struct{}
)
