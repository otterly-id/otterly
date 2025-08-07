package models

type SuccessResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type FailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

type SuccessResponseWithoutData struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}