package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func respondError(w http.ResponseWriter, status int, err error) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(body)
	}
}

func respondNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func parseRequest(r *http.Request, out interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(out)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Failed to parse JSON input: %s", err)
	}

	return err
}
