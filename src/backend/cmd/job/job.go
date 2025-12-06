package job

import (
	"fmt"
	"local/config"
	"local/util/logger"

	"github.com/spf13/cobra"
)

var JobCmd = &cobra.Command{
	Use:   "job",
	Short: "Run background jobs",
	Long:  "Run background jobs for the application",
	Run: func(cmd *cobra.Command, args []string) {
		RunJob()
	},
}

var (
	jobType string
)

func init() {
	JobCmd.Flags().StringVarP(&jobType, "type", "t", "", "Type of job to run (cleanup, sync)")
	JobCmd.MarkFlagRequired("type")
}

func RunJob() {
	// Load config
	config.LoadConfig()

	// Initialize logger
	_, err := logger.InitTracer("job-service")
	if err != nil {
		fmt.Printf("Failed to initialize tracer: %v\n", err)
		return
	}
	defer func() {
		if err := logger.Shutdown(); err != nil {
			fmt.Printf("Error shutting down tracer: %v\n", err)
		}
	}()

	logger.Info(nil, "Job service started", map[string]interface{}{
		"job_type": jobType,
	})

	// Job logic here
	switch jobType {
	case "cleanup":
		runCleanupJob()
	case "sync":
		runSyncJob()
	default:
		fmt.Printf("Unknown job type: %s\n", jobType)
		fmt.Println("Available job types: cleanup, sync")
	}
}

func runCleanupJob() {
	logger.Info(nil, "Running cleanup job")
	// Cleanup job logic
	fmt.Println("Cleanup job executed")
}

func runSyncJob() {
	logger.Info(nil, "Running sync job")
	// Sync job logic
	fmt.Println("Sync job executed")
}

