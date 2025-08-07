package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

func StartServerWithGracefulShutdown(server *http.Server, log *zap.Logger) {
	var wait time.Duration

	go func() {
		log.Info("Server is starting...")
		if err := server.ListenAndServe(); err != nil {
			log.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}

	log.Info("Server is shutting down...")
}
