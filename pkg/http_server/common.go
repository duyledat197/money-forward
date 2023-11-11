package http_server

import (
	"path/filepath"
	"regexp"
)

const (
	slash = string(filepath.Separator)
	space = " "
)

var bracketRegex = regexp.MustCompile(`\{(.*?)\}`)
