package roundrobin

import (
	"log"
	"net/http"
	"pblb/lib"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	asyncHealthChecksTimeSeconds = 15
)

var (
	rrProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pblb_roundrobin_processed_total",
			Help: "The total number of processed requests",
		},
		[]string{"status_class"},
	)
	rrHealthyNodes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pblb_roundrobin_healthy_nodes",
		Help: "The total number of healthy nodes",
	})
	rrTotalNodes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pblb_roundrobin_total_nodes",
		Help: "The total number of healthy and unhealthy nodes",
	})
)

// RoundRobin struct contains:
// - A slice of node pointers
// - The current "active" node index
// - The maximum number of Healthy nodes
type RoundRobin struct {
	Nodes   []*lib.Node
	current int
	max     int
}

// New creates a new RoundRobin load balancer.
func New(nodes []*lib.Node) RoundRobin {
	rr := RoundRobin{nodes, 0, len(nodes)}
	rrHealthyNodes.Set(float64(rr.max))
	rrTotalNodes.Set(float64(rr.max))
	rr.AsyncHealthChecks()

	return rr
}

// AsyncHealthChecks performs health checks in the background at an interval
// set by asyncHealthChecksTimeSeconds.
func (rr *RoundRobin) AsyncHealthChecks() {
	go func() {
		for {
			log.Println("Performing async health checks")
			healthyNodes := 0
			for _, n := range rr.Nodes {
				healthy := n.CheckHealth()
				if healthy {
					healthyNodes++
				}
			}
			log.Printf("%d out of %d nodes are healthy", healthyNodes, len(rr.Nodes))
			rrHealthyNodes.Set(float64(healthyNodes))
			time.Sleep(asyncHealthChecksTimeSeconds * time.Second)
		}
	}()
}

// Handler selects a node via round robin and passes the request to the
// selected node.
func (rr *RoundRobin) Handler(w http.ResponseWriter, r *http.Request) {
	node := rr.selectNode()

	log.Printf("Handling request to %s:%s. Active Connections: %d. Method: RoundRobin.\n", node.Address, node.Port, node.ActiveConnections)

	switch status := node.Handler(w, r); {
	case status >= 500:
		log.Printf("Node %s:%s failed to process request. Status: %d.\n", node.Address, node.Port, status)
		rr.idempotentDeactivateNode(node)
		rrProcessed.WithLabelValues("5xx").Inc()
	case status >= 400:
		rr.idempotentRecoverNode(node)
		rrProcessed.WithLabelValues("4xx").Inc()
	case status >= 300:
		rr.idempotentRecoverNode(node)
		rrProcessed.WithLabelValues("3xx").Inc()
	case status >= 200:
		rr.idempotentRecoverNode(node)
		rrProcessed.WithLabelValues("2xx").Inc()
	default:
		rr.idempotentRecoverNode(node)
		rrProcessed.WithLabelValues("1xx").Inc()
	}
}

func (rr *RoundRobin) selectNode() *lib.Node {
	node := rr.Nodes[rr.current]
	count := 0

	// If there's no healthy nodes, just serve round robin. Otherwise, iterate
	// until we get a healthy node.
	//
	// TODO: We should set an "inversion" value in config.yaml when, if the sum
	// of healthy nodes is equal or below this value, we just serve to all nodes.
	for node.IsUnhealthy() && count < rr.max {
		rr.current = (rr.current + 1) % rr.max
		node = rr.Nodes[rr.current]
		count++
	}

	rr.current = (rr.current + 1) % rr.max
	return node
}

func (rr *RoundRobin) idempotentRecoverNode(n *lib.Node) {
	if n.IsUnhealthy() {
		n.SetHealthy()
		rrHealthyNodes.Inc()
	}
}

func (rr *RoundRobin) idempotentDeactivateNode(n *lib.Node) {
	if n.IsHealthy() {
		n.SetUnhealthy()
		rrHealthyNodes.Dec()
	}
}
