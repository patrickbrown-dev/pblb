package server

import (
	"fmt"
	"log"
	"net/http"
	"pblb/lib"

	"github.com/spf13/viper"
)

// Serve sets up the http server for the database.
func Serve(lb lib.LoadBalancer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lb.Handler(w, r)
	})
	port := fmt.Sprintf(":%s", viper.GetString("port"))
	log.Printf("Starting pblb server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
