package lib

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// RoundRobin struct contains:
// - A slice of node pointers
// - The current "active" node index
// - The maximum number of Healthy nodes
type RoundRobin struct {
	Nodes   []*Node
	current int
	max     int
	mux     sync.Mutex
}

// NewRoundRobin creates a new RoundRobin load balancer.
func NewRoundRobin(nodes []*Node) *RoundRobin {
	rr := RoundRobin{Nodes: nodes, current: 0, max: len(nodes)}
	healthyNodesGauge.Set(float64(rr.max))
	totalNodesGauge.Set(float64(rr.max))
	rr.AsyncHealthChecks()

	return &rr
}

// AsyncHealthChecks performs health checks in the background at an interval
// set by asyncHealthChecksTimeSeconds.
func (rr *RoundRobin) AsyncHealthChecks() {
	go func() {
		for {
			log.Println("Performing async health checks")
			healthyNodes := 0
			for _, n := range rr.Nodes {
				if n.CheckHealth() {
					rr.idempotentRecoverNode(n)
					healthyNodes++
				} else {
					rr.idempotentDeactivateNode(n)
				}
			}
			log.Printf("%d out of %d nodes are healthy", healthyNodes, len(rr.Nodes))
			healthyNodesGauge.Set(float64(healthyNodes))
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
		processedTotal.WithLabelValues("5xx", node.Address).Inc()
	case status >= 400:
		rr.idempotentRecoverNode(node)
		processedTotal.WithLabelValues("4xx", node.Address).Inc()
	case status >= 300:
		rr.idempotentRecoverNode(node)
		processedTotal.WithLabelValues("3xx", node.Address).Inc()
	case status >= 200:
		rr.idempotentRecoverNode(node)
		processedTotal.WithLabelValues("2xx", node.Address).Inc()
	default:
		rr.idempotentRecoverNode(node)
		processedTotal.WithLabelValues("1xx", node.Address).Inc()
	}
}

func (rr *RoundRobin) selectNode() *Node {
	rr.mux.Lock()
	defer rr.mux.Unlock()
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

func (rr *RoundRobin) idempotentRecoverNode(n *Node) {
	if n.IsUnhealthy() {
		n.SetHealthy()
		healthyNodesGauge.Inc()
	}
}

func (rr *RoundRobin) idempotentDeactivateNode(n *Node) {
	if n.IsHealthy() {
		n.SetUnhealthy()
		healthyNodesGauge.Dec()
	}
}
