package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/app/controllers"
)

func PublicRoutes(router chi.Router) {
	router.Route("/users", func(r chi.Router) {
		r.Post("/", controllers.CreateUser)
		r.Get("/", controllers.GetUsers)
		r.Get("/{id}", controllers.GetUser)
		r.Patch("/{id}", controllers.UpdateUser)
		r.Delete("/{id}", controllers.DeleteUser)
	})
}
