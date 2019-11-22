package cmd

import (
	"fmt"
	"log"
	"pblb/lib"
	"pblb/server"
	"pblb/twochoice"

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

	lb := twochoice.New(nodes)

	server.Serve(&lb)
}
