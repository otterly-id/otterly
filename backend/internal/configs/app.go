package configs

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/otterly-id/otterly/backend/db"
	"github.com/otterly-id/otterly/backend/internal/api/controllers"
	"github.com/otterly-id/otterly/backend/internal/delivery/middlewares"
	"github.com/otterly-id/otterly/backend/internal/delivery/route"
	"github.com/otterly-id/otterly/backend/internal/helpers"
	"github.com/otterly-id/otterly/backend/internal/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type BootstrapConfig struct {
	App      chi.Router
	Log      *zap.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Server   *http.Server
	DB       *db.Queries
}

func Bootstrap(config *BootstrapConfig) {
	jwtSecret := config.Config.GetString("JWT_SECRET")
	jwtExpiresIn := config.Config.GetInt("JWT_EXPIRES_IN")

	jwtManager := utils.NewJWTManager(
		[]byte(jwtSecret),
		"otterly-backend",
		"otterly-users",
		time.Duration(jwtExpiresIn)*time.Hour,
	)

	responseHandler := helpers.NewHandler(config.Log)

	userController := controllers.NewUserController(config.Log, config.Validate, config.DB)
	authController := controllers.NewAuthController(config.Log, config.Validate, config.DB, jwtManager)

	authMiddleware := middlewares.NewAuthMiddleware(jwtManager, responseHandler, config.Log)

	routeConfig := route.RouteConfig{
		App:             config.App,
		Log:             config.Log,
		UserController:  userController,
		AuthController:  authController,
		ResponseHandler: helpers.NewHandler(config.Log),
		AuthMiddleware:  authMiddleware,
	}

	routeConfig.Setup()

	utils.StartServerWithGracefulShutdown(config.Server, config.Log)
}