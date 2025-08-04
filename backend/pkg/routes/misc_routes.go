package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/pkg/helpers"
	"go.uber.org/zap"
)

func MiscRoutes(router chi.Router, logger *zap.Logger) {
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Route doesn't exist",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr))

		helpers.NewFailureResponse[any]("Route doesn't exist", nil).Write(
			w, &helpers.ResponseOptions{StatusCode: http.StatusNotFound})
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Route doesn't exist",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr))

		helpers.NewFailureResponse[any]("Method not allowed", nil).Write(
			w, &helpers.ResponseOptions{StatusCode: http.StatusMethodNotAllowed})
	})
}
