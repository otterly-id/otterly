package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Envelope[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("could not encode response to JSON: %v", err)
	}
}

func SuccessResponse[T any](w http.ResponseWriter, statusCode int, message string, data T) {
	response := Envelope[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
	writeJSON(w, statusCode, response)
}

func FailureResponse(w http.ResponseWriter, statusCode int, message string, errors any) {
	response := Envelope[any]{
		Success: false,
		Message: message,
		Errors:  errors,
	}
	writeJSON(w, statusCode, response)
}
