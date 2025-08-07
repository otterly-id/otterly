package configs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/otterly-id/otterly/backend/db"
	"github.com/otterly-id/otterly/backend/internal/api/controllers"
	"github.com/otterly-id/otterly/backend/internal/helpers"
	"github.com/otterly-id/otterly/backend/internal/route"
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
	userController := controllers.NewUserController(config.Log, config.Validate, config.DB)

	routeConfig := route.RouteConfig{
		App:             config.App,
		Log:             config.Log,
		UserController:  userController,
		ResponseHandler: helpers.NewHandler(config.Log),
	}

	routeConfig.Setup()

	utils.StartServerWithGracefulShutdown(config.Server, config.Log)
}
