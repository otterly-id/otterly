package configs

import (
	"net/http"

	"github.com/go-chi/cors"
)

func NewCORS() func(http.Handler) http.Handler {
	corsConfig := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}

	return cors.Handler(corsConfig)
}
