package http

import (
	"encoding/json"
	"net/http"
)

// HTTPServer provides HTTP API endpoints
type HTTPServer struct {
	address string
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(address string) *HTTPServer {
	return &HTTPServer{
		address: address,
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	http.HandleFunc("/health", s.handleHealth)
	http.HandleFunc("/metrics", s.handleMetrics)
	http.HandleFunc("/status", s.handleStatus)

	return http.ListenAndServe(s.address, nil)
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (s *HTTPServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: Return Prometheus metrics
	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_players": 0,
		"active_rooms":   0,
	})
}

func (s *HTTPServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"server": "running",
		"uptime": 0,
	})
}
