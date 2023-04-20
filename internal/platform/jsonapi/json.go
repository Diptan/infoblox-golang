package jsonapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// Read unmarshals the JSON-encoded data into v
func Read(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err := r.Body.Close(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

// Write writes a JSON representation of v to response
func Write(w http.ResponseWriter, v interface{}) {
	content, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
