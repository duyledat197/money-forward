// Package http_server provides an easy implementation handler
// for http server with register handle by method, path, generic handler
package http_server

import (
	"bytes"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_retrieveDataFromRequest(t *testing.T) {

	mockBody := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "Dat",
		Age:  26,
	}

	path := "/users/123?job=senior-software-engineer"
	pattern := "/users/{id}"

	b, err := json.Marshal(&mockBody)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	resp := httptest.NewRecorder()

	params, err := retrieveDataFromRequest(resp, appendWildCardParams(pattern, req))
	require.NoError(t, err)

	expectedParams := map[string]any{
		"id":   "123",
		"name": "Dat",
		// because golang json always parse number to float64
		"age": float64(26),
		"job": "senior-software-engineer",
	}

	require.True(t, maps.Equal(params, expectedParams))
}
