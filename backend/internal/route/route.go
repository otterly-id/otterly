package route

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"

	"github.com/otterly-id/otterly/backend/internal/api/controllers"
	"github.com/otterly-id/otterly/backend/internal/helpers"
	"go.uber.org/zap"
)

type RouteConfig struct {
	App             chi.Router
	Log             *zap.Logger
	ResponseHandler *helpers.ResponseHandler
	UserController  *controllers.UserController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupHealthCheckRoute()
	c.SetupDefaultRoute()
	c.SetupSwaggerRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", c.UserController.CreateUser)
			r.Get("/", c.UserController.GetUsers)
			r.Get("/{id}", c.UserController.GetUser)
			r.Patch("/{id}", c.UserController.UpdateUser)
			r.Delete("/{id}", c.UserController.DeleteUser)
		})
	})
}

func (c *RouteConfig) SetupHealthCheckRoute() {
	c.App.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		c.Log.Info("Health check requested",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr))

		c.ResponseHandler.Success(w, r, http.StatusOK, "Service up and running", time.Now().Format(time.RFC1123))
	})
}

func (c *RouteConfig) SetupDefaultRoute() {
	c.App.NotFound(func(w http.ResponseWriter, r *http.Request) {
		c.Log.Info("Route doesn't exist",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr))

		c.ResponseHandler.CustomError(w, r, http.StatusNotFound, "Route doesn't exist", nil)
	})

	c.App.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		c.Log.Info("Method not allowed",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr))

		c.ResponseHandler.CustomError(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
	})
}

func (c *RouteConfig) SetupSwaggerRoute() {
	c.App.Get("/", func(w http.ResponseWriter, r *http.Request) {
		requestStart := time.Now()

		c.Log.Info("Swagger documentation requested",
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
			c.Log.Error("Failed to generate Swagger HTML",
				zap.Error(err),
				zap.Duration("request_duration", time.Since(requestStart)))

			c.ResponseHandler.CustomError(w, r, http.StatusInternalServerError, "Failed to generate API reference HTML", fmt.Errorf("Unable to load API documentation"))
			return
		}

		c.Log.Info("Swagger documentation served",
			zap.Duration("request_duration", time.Since(requestStart)),
			zap.Int("content_length", len(htmlContent)))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))
	})
}
