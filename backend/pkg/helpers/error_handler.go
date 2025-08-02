package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/otterly-id/otterly/backend/pkg/utils"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (eh *ErrorHandler) JSONDecodeError(w http.ResponseWriter, r *http.Request, err error) {
	eh.logger.Error("JSON decode error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	NewFailureResponse("Failed to parse JSON body", "Invalid JSON format").Write(
		w, &ResponseOptions{StatusCode: http.StatusBadRequest})
}

func (eh *ErrorHandler) ValidationError(w http.ResponseWriter, r *http.Request, err error) {
	eh.logger.Error("Validation error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	validationErrors := utils.ValidatorErrors(err)
	NewFailureResponse("Validation failed", validationErrors).Write(
		w, &ResponseOptions{StatusCode: http.StatusBadRequest})
}

func (eh *ErrorHandler) DBConnectionError(w http.ResponseWriter, r *http.Request, err error) {
	eh.logger.Error("Database connection error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	NewFailureResponse("Service temporarily unavailable", "Please try again later").Write(
		w, nil)
}

func (eh *ErrorHandler) InvalidIDError(w http.ResponseWriter, r *http.Request, err error) {
	eh.logger.Error("Invalid ID error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("id", chi.URLParam(r, "id")),
		zap.Error(err))

	NewFailureResponse("Invalid ID format", "The provided ID is not in the correct format").Write(
		w, &ResponseOptions{StatusCode: http.StatusBadRequest})
}

func (eh *ErrorHandler) NotFoundError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	eh.logger.Info("Resource not found",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("%s not found", resource)
	errorDetail := fmt.Sprintf("The requested %s could not be found", strings.ToLower(resource))

	NewFailureResponse(message, errorDetail).Write(
		w, &ResponseOptions{StatusCode: http.StatusNotFound})
}

func (eh *ErrorHandler) DuplicateKeyError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	eh.logger.Error("Duplicate key error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("%s already exists", strings.Title(resource))
	errorDetail := fmt.Sprintf("A %s with this information already exists", strings.ToLower(resource))

	NewFailureResponse(message, errorDetail).Write(
		w, &ResponseOptions{StatusCode: http.StatusConflict})
}

func (eh *ErrorHandler) PasswordError(w http.ResponseWriter, r *http.Request, err error) {
	eh.logger.Error("Password error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	NewFailureResponse("Authentication failed", "Invalid credentials provided").Write(
		w, &ResponseOptions{StatusCode: http.StatusUnauthorized})
}

func (eh *ErrorHandler) CreateError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
		eh.DuplicateKeyError(w, r, err, resource)
		return
	}

	eh.logger.Error("Create error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("Failed to create %s", strings.ToLower(resource))
	errorDetail := fmt.Sprintf("An error occurred while creating the %s", strings.ToLower(resource))

	NewFailureResponse(message, errorDetail).Write(
		w, &ResponseOptions{StatusCode: http.StatusInternalServerError})
}

func (eh *ErrorHandler) UpdateError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if err == pgx.ErrNoRows {
		eh.NotFoundError(w, r, err, resource)
		return
	}

	eh.logger.Error("Update error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("Failed to update %s", strings.ToLower(resource))
	errorDetail := fmt.Sprintf("An error occurred while updating the %s", strings.ToLower(resource))

	NewFailureResponse(message, errorDetail).Write(
		w, &ResponseOptions{StatusCode: http.StatusInternalServerError})
}

func (eh *ErrorHandler) DeleteError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if err == pgx.ErrNoRows {
		eh.NotFoundError(w, r, err, resource)
		return
	}

	eh.logger.Error("Delete error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("Failed to delete %s", strings.ToLower(resource))
	errorDetail := fmt.Sprintf("An error occurred while deleting the %s", strings.ToLower(resource))

	NewFailureResponse(message, errorDetail).Write(
		w, &ResponseOptions{StatusCode: http.StatusInternalServerError})
}
