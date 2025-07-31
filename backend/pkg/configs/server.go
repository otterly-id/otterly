package configs

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/pkg/utils"
)

func ServerConfig(router chi.Router) *http.Server {
	serverConnURL, _ := utils.ConnectionURLBuilder("server")
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	return &http.Server{
		Handler:     router,
		Addr:        serverConnURL,
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}
