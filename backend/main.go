package main

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/otterly-id/otterly/backend/pkg/configs"
	"github.com/otterly-id/otterly/backend/pkg/routes"
	"github.com/otterly-id/otterly/backend/pkg/utils"
)

// @title           Otterly API
// @version         1.0
// @description     Official Otterly API documentation.
// @BasePath  /
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found")
		return
	}

	logger := utils.NewLogger()

	defer logger.Sync()

	router := chi.NewRouter()

	router.Use(cors.Handler(configs.CORSConfig()))

	routes.HealthCheckRoute(router, logger)
	routes.SwaggerRoute(router, logger)
	routes.MiscRoutes(router, logger)

	router.Route("/api", func(r chi.Router) {
		routes.PublicRoutes(r, logger)
	})

	server := configs.ServerConfig(router)

	utils.StartServerWithGracefulShutdown(server)
}
