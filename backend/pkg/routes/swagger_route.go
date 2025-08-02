package routes

import (
	"net/http"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/pkg/helpers"
	"go.uber.org/zap"
)

func SwaggerRoute(router chi.Router, logger *zap.Logger) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		requestStart := time.Now()

		logger.Info("Swagger documentation requested",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr))

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Otterly API",
			},
			DarkMode: true,
			Theme:    "moon",
		})

		if err != nil {
			logger.Error("Failed to generate Swagger HTML",
				zap.Error(err),
				zap.Duration("request_duration", time.Since(requestStart)))

			helpers.NewFailureResponse("Failed to generate API reference HTML", "Unable to load API documentation").Write(
				w, &helpers.ResponseOptions{StatusCode: http.StatusInternalServerError})
			return
		}

		logger.Info("Swagger documentation served",
			zap.Duration("request_duration", time.Since(requestStart)),
			zap.Int("content_length", len(htmlContent)))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))
	})
}
