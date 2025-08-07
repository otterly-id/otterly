package main

import (
	"github.com/otterly-id/otterly/backend/db"
	"github.com/otterly-id/otterly/backend/internal/configs"
	"go.uber.org/zap"
)

// @title           Otterly API
// @version         1.0
// @description     Official Otterly API documentation.
// @BasePath  /
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
