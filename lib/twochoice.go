package lib

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

// TwoChoice struct contains:
// - A slice of node pointers
// - A set of indexes to healthy nodes (as a map)
// - A set of indexes to unhealthy ndoes (as a map)
type TwoChoice struct {
	Nodes          []*Node
	HealthyNodes   map[int]bool
	UnhealthyNodes map[int]bool
}

// NewTwoChoice creates a new TwoChoice load balancer
func NewTwoChoice(nodes []*Node) *TwoChoice {
	if len(nodes) < 3 {
		log.Fatalf("At least 3 nodes are required for TwoChoice, %d found", len(nodes))
	}

	tc := TwoChoice{nodes, make(map[int]bool), make(map[int]bool)}
	rand.Seed(time.Now().UTC().UnixNano())

	for i, n := range tc.Nodes {
		n.SetHealthy()
		tc.HealthyNodes[i] = true
	}

	healthyNodesGauge.Set(float64(len(nodes)))
	totalNodesGauge.Set(float64(len(nodes)))
	tc.AsyncHealthChecks()

	return &tc
}

// AsyncHealthChecks performs health checks in the background at an interval
// set by asyncHealthChecksTimeSeconds.
func (tc *TwoChoice) AsyncHealthChecks() {
	go func() {
		for {
			log.Println("Performing async health checks")
			tc.healthChecks()
			time.Sleep(asyncHealthChecksTimeSeconds * time.Second)
		}
	}()
}

func (tc *TwoChoice) healthChecks() {
	for i, n := range tc.Nodes {
		if n.CheckHealth() {
			tc.idempotentRecoverNode(n, i)
		} else {
			tc.idempotentDeactivateNode(n, i)
		}
	}
	log.Printf("%d out of %d nodes are healthy", len(tc.HealthyNodes), len(tc.Nodes))
}

// Handler selects a node via random two choice and passes the request to the
// selected node. See https://www.nginx.com/blog/nginx-power-of-two-choices-load-balancing-algorithm/
func (tc *TwoChoice) Handler(w http.ResponseWriter, r *http.Request) {
	nodeKey := tc.selectNodeKey()
	node := tc.Nodes[nodeKey]

	log.Printf("Handling request to %s:%s. Active Connections: %d. Method: TwoChoice.\n", node.Address, node.Port, node.ActiveConnections)

	switch status := node.Handler(w, r); {
	case status >= 500:
		log.Printf("Node %s:%s failed to process request. Status: %d.\n", node.Address, node.Port, status)
		tc.idempotentDeactivateNode(node, nodeKey)
		processedTotal.WithLabelValues("5xx", node.Address).Inc()
	case status >= 400:
		tc.idempotentRecoverNode(node, nodeKey)
		processedTotal.WithLabelValues("4xx", node.Address).Inc()
	case status >= 300:
		tc.idempotentRecoverNode(node, nodeKey)
		processedTotal.WithLabelValues("3xx", node.Address).Inc()
	case status >= 200:
		tc.idempotentRecoverNode(node, nodeKey)
		processedTotal.WithLabelValues("2xx", node.Address).Inc()
	default:
		tc.idempotentRecoverNode(node, nodeKey)
		processedTotal.WithLabelValues("1xx", node.Address).Inc()
	}
}

func (tc *TwoChoice) selectNodeKey() int {
	var nodePool map[int]bool

	// If we have less than 2 healthy nodes, serve to the unhealthy node pool.
	if len(tc.HealthyNodes) >= 2 {
		nodePool = tc.HealthyNodes
	} else {
		nodePool = tc.UnhealthyNodes
	}

	keys := make([]int, len(nodePool))

	i := 0
	for k := range nodePool {
		keys[i] = k
		i++
	}

	first := keys[rand.Intn(len(keys))]
	second := keys[rand.Intn(len(keys))]

	for first == second {
		second = keys[rand.Intn(len(keys))]
	}

	log.Printf(
		"TwoChoice Candidates: %s:%s (ActiveConnections: %d), %s:%s (ActiveConnections: %d)",
		tc.Nodes[first].Address,
		tc.Nodes[first].Port,
		tc.Nodes[first].ActiveConnections,
		tc.Nodes[second].Address,
		tc.Nodes[second].Port,
		tc.Nodes[second].ActiveConnections,
	)

	if tc.Nodes[first].ActiveConnections < tc.Nodes[second].ActiveConnections {
		return first
	}
	return second
}

func (tc *TwoChoice) idempotentRecoverNode(n *Node, key int) {
	if n.IsUnhealthy() {
		tc.HealthyNodes[key] = true
		delete(tc.UnhealthyNodes, key)
		n.SetHealthy()
		healthyNodesGauge.Inc()
	}
}

func (tc *TwoChoice) idempotentDeactivateNode(n *Node, key int) {
	if n.IsHealthy() {
		tc.UnhealthyNodes[key] = true
		delete(tc.HealthyNodes, key)
		n.SetUnhealthy()
		healthyNodesGauge.Dec()
	}
}
