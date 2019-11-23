package cmd

import (
	"fmt"
	"log"
	"net/http"
	"pblb/lib"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the pblb server",
	Run:   run,
}

func init() {
	viper.SetDefault("port", "2839")

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/pblb/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func run(cmd *cobra.Command, args []string) {
	var nodes []*lib.Node
	err := viper.UnmarshalKey("nodes", &nodes)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	for _, n := range nodes {
		n.Init()
	}

	switch method := viper.GetString("method"); {
	case method == "roundrobin":
		log.Println("Using the RoundRobin load balancing method")
		lb := lib.NewRoundRobin(nodes)
		serve(&lb)
	case method == "twochoice":
		log.Println("Using the TwoChoice load balancing method")
		lb := lib.NewTwoChoice(nodes)
		serve(&lb)
	default:
		log.Fatalf("Could not find a matching load balancing method to configuration \"%s\"", method)
	}
}

func serve(lb lib.LoadBalancer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lb.Handler(w, r)
	})
	http.Handle("/metrics", promhttp.Handler())

	port := fmt.Sprintf(":%s", viper.GetString("port"))
	log.Printf("Starting pblb server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
