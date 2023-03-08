package mgmt

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ListenAndServe starts the management server to expose metrics, config and
// health status
func ListenAndServe(addr string) error {
	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, router)
}
