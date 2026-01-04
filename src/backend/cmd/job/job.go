package job

import (
	"fmt"
	"local/config"
	"local/job/consumer"
	"local/job/scheduler"
	"local/job/worker"
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
	JobCmd.Flags().StringVarP(&jobType, "type", "t", "", "Type of job to run (cleanup, sync, temporal, kafka-message, kafka-notification)")
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
	case "temporal":
		runTemporalWorker()
	case "kafka-message":
		runKafkaMessageConsumer()
	case "kafka-notification":
		runKafkaNotificationConsumer()
	default:
		fmt.Printf("Unknown job type: %s\n", jobType)
		fmt.Println("Available job types: cleanup, sync, temporal, kafka-message, kafka-notification")
	}
}

func runCleanupJob() {
	scheduler.RunCleanupJob()
}

func runSyncJob() {
	scheduler.RunSyncJob()
}

func runTemporalWorker() {
	logger.Info(nil, "Starting Temporal worker service", nil)

	// Create and start Temporal worker
	workerService, err := worker.NewWorkerService()
	if err != nil {
		logger.Error(nil, "Failed to create Temporal worker", err)
		fmt.Printf("Failed to create Temporal worker: %v\n", err)
		return
	}

	// Ensure cleanup on exit
	defer workerService.Stop()

	// Start worker (blocks until interrupted)
	if err := workerService.Start(); err != nil {
		logger.Error(nil, "Temporal worker error", err)
		fmt.Printf("Temporal worker error: %v\n", err)
		return
	}

	logger.Info(nil, "Temporal worker service stopped", nil)
}

func runKafkaMessageConsumer() {
	logger.Info(nil, "Starting Kafka message consumer", nil)

	if err := consumer.StartMessageConsumer(); err != nil {
		logger.Error(nil, "Kafka message consumer error", err)
		fmt.Printf("Kafka message consumer error: %v\n", err)
		return
	}

	logger.Info(nil, "Kafka message consumer stopped", nil)
}

func runKafkaNotificationConsumer() {
	logger.Info(nil, "Starting Kafka notification consumer", nil)

	if err := consumer.StartNotificationConsumer(); err != nil {
		logger.Error(nil, "Kafka notification consumer error", err)
		fmt.Printf("Kafka notification consumer error: %v\n", err)
		return
	}

	logger.Info(nil, "Kafka notification consumer stopped", nil)
}

