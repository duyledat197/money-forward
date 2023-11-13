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
	userInfoKey       struct{}

	UserInfo struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
)
