package routes

import (
	"math"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/otterly-id/otterly/backend/pkg/utils"
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
	utils.SuccessResponse
	Data MemStats `json:"data"`
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func byteToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// @Description   Check API health
// @Produce       json
// @Success       200   {object} WrappedResponse
// @Tags		  Health Check
// @Router        /health-check [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
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

	utils.NewSuccessResponse("Service up and running!", *memstats).Write(w, nil)
}

func HealthCheckRoute(router chi.Router) {
	router.Get("/health-check", healthCheck)
}
