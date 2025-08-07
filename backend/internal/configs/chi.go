package configs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewChi(cors func(http.Handler) http.Handler) chi.Router {
	c := chi.NewRouter()
	c.Use(cors)
	return c
}
