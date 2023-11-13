package http_server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"user-management/pkg/token_utils"
)

// appendWildCardParams will get wildcard params in request and append to the request context
func appendWildCardParams(pattern string, r *http.Request) *http.Request {
	result := make(map[string]any)
	patternEls := strings.Split(strings.Trim(pattern, slash), slash)
	sourceEls := strings.Split(strings.Trim(r.URL.Path, slash), slash)
	for i, el := range patternEls {
		if bracketRegex.MatchString(el) {
			val := strings.TrimLeft(el, openBracket)
			val = strings.TrimRight(val, closeBracket)
			result[val] = sourceEls[i]
		}
	}

	return r.WithContext(context.WithValue(r.Context(), &wildcardParamsKey{}, result))
}

// ImportUserInfoToContext implements import the user info which retrieved from token
// and inject it into the given context.
func ImportUserInfoToContext(ctx context.Context, info *token_utils.Payload) context.Context {
	return context.WithValue(ctx, &userInfoKey{}, info)
}

// ExtractUserInfoFromContext returns an user info which was injected from [ImportUserInfoToContext].
func ExtractUserInfoFromContext(ctx context.Context) (*token_utils.Payload, error) {
	info, ok := ctx.Value(&userInfoKey{}).(*token_utils.Payload)

	if !ok || info == nil {
		return nil, fmt.Errorf("authorization is not valid")
	}

	return info, nil
}
