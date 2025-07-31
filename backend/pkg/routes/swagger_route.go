package routes

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/pkg/utils"
)

func SwaggerRoute(router chi.Router) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Otterly API",
			},
			DarkMode: true,
		})

		if err != nil {
			utils.NewFailureResponse(
				"Failed to generate API reference HTML", err.Error()).Write(
				w, &utils.ResponseOptions{StatusCode: http.StatusInternalServerError})

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlContent))
	})
}
