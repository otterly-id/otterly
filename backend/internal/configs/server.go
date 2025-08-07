package configs

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func NewServer(config *viper.Viper, router chi.Router) *http.Server {
	serverConnURL := config.GetString("SERVER_URL")
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	return &http.Server{
		Handler:     router,
		Addr:        serverConnURL,
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}
