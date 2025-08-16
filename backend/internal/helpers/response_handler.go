package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/otterly-id/otterly/backend/internal/utils"
	"go.uber.org/zap"
)

type ResponseHandler struct {
	Log *zap.Logger
}

func NewHandler(log *zap.Logger) *ResponseHandler {
	return &ResponseHandler{
		Log: log,
	}
}

func (rh *ResponseHandler) Success(w http.ResponseWriter, r *http.Request, statusCode int, message string, data any) {
	rh.Log.Info(message,
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method))
	utils.SuccessResponse(w, statusCode, message, data)
}

func (rh *ResponseHandler) JSONDecodeError(w http.ResponseWriter, r *http.Request, err error) {
	rh.Log.Error("JSON decode error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	utils.FailureResponse(w, http.StatusBadRequest, "Failed to parse JSON body", "Invalid JSON format")
}

func (rh *ResponseHandler) ValidationError(w http.ResponseWriter, r *http.Request, err error) {
	rh.Log.Error("Validation error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	validationErrors := ValidatorErrors(err)
	utils.FailureResponse(w, http.StatusBadRequest, "Validation failed", validationErrors)
}

func (rh *ResponseHandler) InvalidIDError(w http.ResponseWriter, r *http.Request, err error) {
	rh.Log.Error("Invalid ID error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("id", chi.URLParam(r, "id")),
		zap.Error(err))

	utils.FailureResponse(w, http.StatusBadRequest, "Invalid ID format", "The provided ID is not in the correct format")
}

func (rh *ResponseHandler) NotFoundError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	rh.Log.Info("Resource not found",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("%s not found", resource)
	errorDetail := fmt.Sprintf("The requested %s could not be found", strings.ToLower(resource))

	utils.FailureResponse(w, http.StatusNotFound, message, errorDetail)
}

func (rh *ResponseHandler) DuplicateKeyError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	rh.Log.Error("Duplicate key error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("%s already exists", strings.Title(resource))
	errorDetail := fmt.Sprintf("A %s with this information already exists", strings.ToLower(resource))

	utils.FailureResponse(w, http.StatusConflict, message, errorDetail)
}

func (rh *ResponseHandler) JWTError(w http.ResponseWriter, r *http.Request, err error) {
	rh.Log.Error("JWT error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	utils.FailureResponse(w, http.StatusUnauthorized, "Authentication failed", "Invalid or expired token")
}

func (rh *ResponseHandler) HashPasswordError(w http.ResponseWriter, r *http.Request, err error) {
	rh.Log.Error("Failed to hash password",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))
	utils.FailureResponse(w, http.StatusInternalServerError, "Failed to hash password", "An error occurred while hashing the password")
}

func (rh *ResponseHandler) AuthenticationRequiredError(w http.ResponseWriter, r *http.Request) {
	rh.Log.Error("Authentication required",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method))
	utils.FailureResponse(w, http.StatusUnauthorized, "Authentication required", "You must be authenticated to access this resource")
}

func (rh *ResponseHandler) AuthenticationFailedError(w http.ResponseWriter, r *http.Request, err error) {
	rh.Log.Error("Authentication failed",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))
	utils.FailureResponse(w, http.StatusUnauthorized, "Authentication failed", "Invalid credentials provided")
}

func (rh *ResponseHandler) TokenGenerationError(w http.ResponseWriter, r *http.Request, err error) {
	rh.Log.Error("Failed to generate token",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))
	utils.FailureResponse(w, http.StatusInternalServerError, "Failed to generate token", "An error occurred while generating the authentication token")
}

func (rh *ResponseHandler) InsufficientPermissionsError(w http.ResponseWriter, r *http.Request) {
	rh.Log.Error("Insufficient permissions",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method))
	utils.FailureResponse(w, http.StatusForbidden, "Insufficient permissions", "You do not have permission to access this resource")
}

func (rh *ResponseHandler) CreateItemError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
		rh.DuplicateKeyError(w, r, err, resource)
		return
	}

	rh.Log.Error("Create error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("Failed to create %s", strings.ToLower(resource))
	errorDetail := fmt.Sprintf("An error occurred while creating the %s", strings.ToLower(resource))

	utils.FailureResponse(w, http.StatusInternalServerError, message, errorDetail)
}

func (rh *ResponseHandler) UpdateItemError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if err == pgx.ErrNoRows {
		rh.NotFoundError(w, r, err, resource)
		return
	}

	rh.Log.Error("Update error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("Failed to update %s", strings.ToLower(resource))
	errorDetail := fmt.Sprintf("An error occurred while updating the %s", strings.ToLower(resource))

	utils.FailureResponse(w, http.StatusInternalServerError, message, errorDetail)
}

func (rh *ResponseHandler) DeleteItemError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if err == pgx.ErrNoRows {
		rh.NotFoundError(w, r, err, resource)
		return
	}

	rh.Log.Error("Delete error",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.String("resource", resource),
		zap.Error(err))

	message := fmt.Sprintf("Failed to delete %s", strings.ToLower(resource))
	errorDetail := fmt.Sprintf("An error occurred while deleting the %s", strings.ToLower(resource))

	utils.FailureResponse(w, http.StatusInternalServerError, message, errorDetail)
}

func (rh *ResponseHandler) CustomError(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {
	rh.Log.Error(message,
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Error(err))

	utils.FailureResponse(w, statusCode, message, err)
}