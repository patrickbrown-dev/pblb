package lib

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	asyncHealthChecksTimeSeconds = 15
)

var (
	processedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pblb_processed_total",
			Help: "The total number of processed requests",
		},
		[]string{"status_class", "node"},
	)
	healthyNodesGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pblb_healthy_nodes",
		Help: "The total number of healthy nodes",
	})
	totalNodesGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pblb_total_nodes",
		Help: "The total number of healthy and unhealthy nodes",
	})
)

// LoadBalancer common inteface
type LoadBalancer interface {
	Handler(w http.ResponseWriter, r *http.Request)
	AsyncHealthChecks()
}
