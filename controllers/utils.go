package controllers

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

type ApplicationError struct {
	Error Error `json:"error"`
}

// JsonResponse generates a new HTTP Response in Json format,
// acording to the HTTP status code it receives
func JsonResponse(w http.ResponseWriter, data []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func NewApplicationError(err error) ApplicationError {
	return ApplicationError{
		Error: Error{Message: err.Error()},
	}
}

func ApplicationErrorResponse(rw http.ResponseWriter, err error, statusCode int) {
	data := NewApplicationError(err)
	resp, _ := json.Marshal(data)
	JsonResponse(rw, resp, statusCode)
}
