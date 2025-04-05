package profile

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/widrin/knights/logger"
)

var (
	httpOnce   sync.Once
	httpServer *http.Server
)

func StartHTTPServer(port int) {
	httpOnce.Do(func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}

		go func() {
			logger.Info("Starting metrics server on port %d", port)
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("Metrics server failed: %v", err)
			}
		}()
	})
}
