package lib

import "net/http"

// LoadBalancer common inteface
type LoadBalancer interface {
	Handler(w http.ResponseWriter, r *http.Request)
}
