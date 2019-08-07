package lib

import "net/http"

// LoadBalancer ...
type LoadBalancer interface {
	Handler(w http.ResponseWriter, r *http.Request)
}
