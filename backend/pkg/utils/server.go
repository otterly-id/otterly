package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServerWithGracefulShutdown(server *http.Server) {
	var wait time.Duration

	logger := NewLogger()

	go func() {
		logger.Info("Server is starting...")
		if err := server.ListenAndServe(); err != nil {
			logger.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Server is shutting down...")
}
