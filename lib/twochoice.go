package lib

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	asyncHealthChecksTimeSeconds = 15
)

var (
	tcProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pblb_twochoice_processed_total",
			Help: "The total number of processed requests",
		},
		[]string{"status_class"},
	)
	tcHealthyNodes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pblb_twochoice_healthy_nodes",
		Help: "The total number of healthy nodes",
	})
	tcTotalNodes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pblb_twochoice_total_nodes",
		Help: "The total number of healthy and unhealthy nodes",
	})
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
func NewTwoChoice(nodes []*Node) TwoChoice {
	if len(nodes) < 3 {
		log.Fatalf("At least 3 nodes are required for TwoChoice, %d found", len(nodes))
	}

	tc := TwoChoice{nodes, make(map[int]bool), make(map[int]bool)}
	rand.Seed(time.Now().UTC().UnixNano())

	for i, n := range tc.Nodes {
		n.SetHealthy()
		tc.HealthyNodes[i] = true
	}

	tcHealthyNodes.Set(float64(len(nodes)))
	tcTotalNodes.Set(float64(len(nodes)))
	tc.AsyncHealthChecks()

	return tc
}

// AsyncHealthChecks performs health checks in the background at an interval
// set by asyncHealthChecksTimeSeconds.
func (tc *TwoChoice) AsyncHealthChecks() {
	go func() {
		for {
			log.Println("Performing async health checks")
			healthyNodes := 0
			for i, n := range tc.Nodes {
				healthy := n.CheckHealth()
				if healthy {
					tc.HealthyNodes[i] = true
					delete(tc.UnhealthyNodes, i)
					healthyNodes++
				} else {
					tc.UnhealthyNodes[i] = true
					delete(tc.HealthyNodes, i)
				}
			}
			log.Printf("%d out of %d nodes are healthy", healthyNodes, len(tc.Nodes))
			tcHealthyNodes.Set(float64(healthyNodes))
			time.Sleep(asyncHealthChecksTimeSeconds * time.Second)
		}
	}()
}

// Handler selects a node via random two choice and passes the request to the
// selected node. See https://www.nginx.com/blog/nginx-power-of-two-choices-load-balancing-algorithm/
func (tc *TwoChoice) Handler(w http.ResponseWriter, r *http.Request) {
	node := tc.selectNode()

	log.Printf("Handling request to %s:%s. Active Connections: %d. Method: TwoChoice.\n", node.Address, node.Port, node.ActiveConnections)

	switch status := node.Handler(w, r); {
	case status >= 500:
		log.Printf("Node %s:%s failed to process request. Status: %d.\n", node.Address, node.Port, status)
		tc.idempotentDeactivateNode(node)
		tcProcessed.WithLabelValues("5xx").Inc()
	case status >= 400:
		tc.idempotentRecoverNode(node)
		tcProcessed.WithLabelValues("4xx").Inc()
	case status >= 300:
		tc.idempotentRecoverNode(node)
		tcProcessed.WithLabelValues("3xx").Inc()
	case status >= 200:
		tc.idempotentRecoverNode(node)
		tcProcessed.WithLabelValues("2xx").Inc()
	default:
		tc.idempotentRecoverNode(node)
		tcProcessed.WithLabelValues("1xx").Inc()
	}
}

func (tc *TwoChoice) selectNode() *Node {
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

	var node *Node
	node1 := tc.Nodes[first]
	node2 := tc.Nodes[second]

	log.Printf("TwoChoice Candidates: %s:%s (ActiveConnections: %d), %s:%s (ActiveConnections: %d)", node1.Address, node1.Port, node1.ActiveConnections, node2.Address, node2.Port, node2.ActiveConnections)

	if node1.ActiveConnections < node2.ActiveConnections {
		node = node1
	} else {
		node = node2
	}

	log.Printf("TwoChoice chose %s:%s", node.Address, node.Port)

	return node
}

func (tc *TwoChoice) idempotentRecoverNode(n *Node) {
	if n.IsUnhealthy() {
		n.SetHealthy()
		tcHealthyNodes.Inc()
	}
}

func (tc *TwoChoice) idempotentDeactivateNode(n *Node) {
	if n.IsHealthy() {
		n.SetUnhealthy()
		tcHealthyNodes.Dec()
	}
}
