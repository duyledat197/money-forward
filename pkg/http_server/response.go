package http_server

import (
	"encoding/json"
	"log"
	"net/http"
)

// response struct present response format to http client.
type response struct {
	Code    int      `json:"code"`
	Message string   `json:"message,omitempty"`
	Details []string `json:"details,omitempty"`
	Data    any      `json:"data,omitempty"`
}

// errorResponse write error to http response with passing code and error.
func errorResponse(w http.ResponseWriter, code int, err error) {
	resp := &response{
		Code:    code,
		Message: err.Error(),
	}

	jData, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println("code", code)
	w.WriteHeader(code)
	if _, err := w.Write(jData); err != nil {
		log.Println(err)
	}
}

// dataResponse write response data to http response with passing data.
// The response status code is fixed to 200.
func dataResponse(w http.ResponseWriter, data any) {
	resp := &response{
		Data: data,
	}

	jData, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jData); err != nil {
		log.Println(err)
	}
}
