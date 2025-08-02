package helpers

import (
	"net/http"

	"go.uber.org/zap"
)

type SuccessHandler struct {
	logger *zap.Logger
}

func NewSuccessHandler(logger *zap.Logger) *SuccessHandler {
	return &SuccessHandler{
		logger: logger,
	}
}

func (sh *SuccessHandler) WithData(w http.ResponseWriter, r *http.Request, message string, data any, statusCode *ResponseOptions) {
	sh.logger.Info(message,
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method))
	NewSuccessResponse(message, data).Write(w, statusCode)
}

func (sh *SuccessHandler) WithoutData(w http.ResponseWriter, r *http.Request, message string, statusCode *ResponseOptions) {
	sh.logger.Info(message,
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method))
	NewSuccessResponse[any](message, nil).Write(w, statusCode)
}
