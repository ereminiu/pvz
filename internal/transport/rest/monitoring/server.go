package monitoring

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(port int) *http.Server {
	http.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Handler:      promhttp.Handler(),
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
