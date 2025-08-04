package routes

import (
	"math"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/pkg/helpers"
	"go.uber.org/zap"
)

type MemStats struct {
	Alloc   float64 `json:"alloc"`
	Sys     float64 `json:"sys"`
	LastGc  uint64  `json:"gc_next"`
	NextGc  uint64  `json:"gc_last"`
	CountGc uint64  `json:"gc_cycle"`
	Date    string  `json:"date"`
}

type WrappedResponse struct {
	helpers.SuccessResponse[MemStats]
	Data MemStats `json:"data"`
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func byteToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// @Summary       Health Check
// @Description   Check API health
// @Tags          Health Check
// @Accept        json
// @Produce       json
// @Success       200   {object} WrappedResponse
// @Router        /health-check [get]
func healthCheckHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check requested",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr))

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		memstats := &MemStats{
			Alloc:   roundFloat(byteToMb(m.Alloc), 5),
			Sys:     roundFloat(byteToMb(m.Sys), 0),
			LastGc:  m.LastGC,
			NextGc:  m.NextGC,
			CountGc: uint64(m.NumGC),
			Date:    time.Now().Format(time.RFC1123),
		}

		logger.Info("Health check completed",
			zap.Float64("memory_alloc_mb", memstats.Alloc),
			zap.Float64("memory_sys_mb", memstats.Sys),
			zap.Uint64("gc_count", memstats.CountGc))

		helpers.NewSuccessResponse("Service up and running!", *memstats).Write(w, nil)
	}
}

func HealthCheckRoute(router chi.Router, logger *zap.Logger) {
	router.Get("/health-check", healthCheckHandler(logger))
}
