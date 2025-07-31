package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/pkg/utils"
)

func MiscRoutes(router chi.Router) {
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.NewFailureResponse[interface{}]("Route doesn't exist", nil).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusNotFound})
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.NewFailureResponse[interface{}]("Method not allowed", nil).Write(
			w, &utils.ResponseOptions{StatusCode: http.StatusNotFound})
	})
}
