package mux

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	error
	Code int
}

// Error renders error (if present) to ResponseWriter
func Error(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}
	if err, ok := err.(HTTPError); ok {
		w.WriteHeader(err.Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, _ = w.Write([]byte(err.Error()))
}

// JSON renders object to ResponseWriter as JSON
func JSON(v interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// Plain renders text to ResponseWriter
func Plain(text string, w http.ResponseWriter) error {
	_, err := w.Write([]byte(text))
	return err
}
