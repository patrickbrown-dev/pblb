package roundrobin

import (
	"net/http"
	"pblb/lib"
)

type RoundRobin struct {
	Nodes   []lib.Node
	current int
	max     int
}

func New(nodes []lib.Node) RoundRobin {
	rr := RoundRobin{nodes, 0, len(nodes)}
	return rr
}

func (rr *RoundRobin) Handler(w http.ResponseWriter, r *http.Request) {
	node := rr.Nodes[rr.current]
	rr.current = (rr.current + 1) % rr.max
	node.Handler(w, r)
}
