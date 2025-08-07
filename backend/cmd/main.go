package main

import (
	"github.com/otterly-id/otterly/backend/db"
	"github.com/otterly-id/otterly/backend/internal/configs"
	"go.uber.org/zap"
)

// @title           Otterly API
// @version         1.0
// @description     Official Otterly API documentation.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name otterly_token
// @description JWT token stored in httpOnly cookie for authentication

// @tag.name Auth
// @tag.description Authentication operations

// @tag.name Users
// @tag.description User management operations

// @tag.name Admin
// @tag.description Admin-only operations

// @tag.name Owner
// @tag.description Owner-only operations

// @tag.name Management
// @tag.description Management operations (Admin or Owner)
func main() {
	log := configs.NewLogger()
	defer log.Sync()

	viperConfig := configs.NewViper()
	validate := configs.NewValidator()
	cors := configs.NewCORS()
	app := configs.NewChi(cors)
	server := configs.NewServer(viperConfig, app)

	db, err := db.GetDBConnection(viperConfig)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	configs.Bootstrap(&configs.BootstrapConfig{
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Server:   server,
		DB:       db,
	})
}
