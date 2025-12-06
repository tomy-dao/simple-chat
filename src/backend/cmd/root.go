package cmd

import (
	"local/cmd/job"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "simple-chat",
	Short: "Simple Chat Application",
	Long:  "A simple chat application with HTTP API and background jobs",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(job.JobCmd)
}

