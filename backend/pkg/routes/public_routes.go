package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/app/controllers"
	"go.uber.org/zap"
)

func PublicRoutes(router chi.Router, logger *zap.Logger) {
	var userController = controllers.NewUserController(logger)

	router.Route("/users", func(r chi.Router) {
		r.Post("/", userController.CreateUser)
		r.Get("/", userController.GetUsers)
		r.Get("/{id}", userController.GetUser)
		r.Patch("/{id}", userController.UpdateUser)
		r.Delete("/{id}", userController.DeleteUser)
	})
}
