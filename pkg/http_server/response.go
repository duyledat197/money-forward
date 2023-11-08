package http_server

import (
	"encoding/json"
	"net/http"
)

// response struct present response format to http client
type response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
	Data    any      `json:"data"`
}

// errorResponse write error to http response with passing code and error
func errorResponse(w http.ResponseWriter, code int, err error) {
	resp := &response{
		Code:    code,
		Message: err.Error(),
		Details: []string{},
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

// dataResponse write response data to http response with passing data.
// The response status code fixed to 200
func dataResponse(w http.ResponseWriter, data any) {
	resp := &response{
		Data: data,
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
