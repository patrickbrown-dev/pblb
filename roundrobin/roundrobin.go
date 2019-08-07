package roundrobin

import (
	"log"
	"net/http"
	"pblb/lib"
)

type RoundRobin struct {
	Nodes   []*lib.Node
	current int
	max     int
}

func New(nodes []*lib.Node) RoundRobin {
	// Set all nodes to healthy.
	// TODO perform health check.
	for _, n := range nodes {
		n.SetHealthy()
	}
	rr := RoundRobin{nodes, 0, len(nodes)}

	return rr
}

func (rr *RoundRobin) Handler(w http.ResponseWriter, r *http.Request) {
	node := rr.selectNode()
	log.Printf("Handling request to %s:%s. Method: RoundRobin.\n", node.Address, node.Port)
	node.Handler(w, r)
}

func (rr *RoundRobin) selectNode() *lib.Node {
	node := rr.Nodes[rr.current]
	count := 0

	// If there's no healthy nodes, just serve round robin. Otherwise, iterate
	// until we get a healthy node.
	for node.IsUnhealthy() && count < rr.max {
		rr.current = (rr.current + 1) % rr.max
		node = rr.Nodes[rr.current]
		count++
	}

	rr.current = (rr.current + 1) % rr.max
	return node
}
