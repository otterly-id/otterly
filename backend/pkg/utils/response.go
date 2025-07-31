package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type FailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

type ResponseOptions struct {
	StatusCode int
}

func NewSuccessResponse(message string, data any) *SuccessResponse {
	return &SuccessResponse{Success: true, Message: message, Data: data}
}

func NewFailureResponse(message string, errors any) *FailureResponse {
	return &FailureResponse{Success: false, Message: message, Errors: errors}
}

func (s *SuccessResponse) Write(w http.ResponseWriter, opts *ResponseOptions) {
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

func (f *FailureResponse) Write(w http.ResponseWriter, opts *ResponseOptions) {
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
