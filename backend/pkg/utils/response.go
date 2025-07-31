package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SuccessResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type SuccessResponseWithoutData struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type FailureResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  T      `json:"errors,omitempty"`
}

type ResponseOptions struct {
	StatusCode int
}

func NewSuccessResponse[T any](message string, data T) *SuccessResponse[T] {
	return &SuccessResponse[T]{Success: true, Message: message, Data: data}
}

func NewFailureResponse[T any](message string, errors T) *FailureResponse[T] {
	return &FailureResponse[T]{Success: false, Message: message, Errors: errors}
}

func (s *SuccessResponse[T]) Write(w http.ResponseWriter, opts *ResponseOptions) {
	var statusCode = 200

	if opts != nil {
		statusCode = opts.StatusCode
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(s)

	if err != nil {
		fmt.Fprintf(w, "{}")
	}
}

func (f *FailureResponse[T]) Write(w http.ResponseWriter, opts *ResponseOptions) {
	var statusCode = 500

	if opts != nil {
		statusCode = opts.StatusCode
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(f)

	if err != nil {
		fmt.Fprintf(w, "{}")
	}
}
