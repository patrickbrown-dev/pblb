package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pblb",
	Short: "A very simple experimental load balancer",
	Long:  `A very simple experimental load balancer`,
}

// Execute evaluates the command line arguments and maps them to
// commands and relevant flags.
func Execute() {
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
