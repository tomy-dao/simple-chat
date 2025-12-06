package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Long:  "Start the HTTP API server",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

