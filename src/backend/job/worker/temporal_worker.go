package worker

import (
	"fmt"
	"local/config"
	"local/job/activities"
	"local/job/workflows"
	"local/util/logger"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	// TaskQueue is the name of the Temporal task queue
	TaskQueue = "simple-chat-task-queue"
)

// WorkerService manages the Temporal worker lifecycle
type WorkerService struct {
	client client.Client
	worker worker.Worker
}

// NewWorkerService creates a new Temporal worker service
func NewWorkerService() (*WorkerService, error) {
	// Get Temporal server address from config or env
	temporalAddr := config.Config.TemporalAddress
	if temporalAddr == "" {
		temporalAddr = os.Getenv("TEMPORAL_ADDRESS")
		if temporalAddr == "" {
			temporalAddr = "localhost:7233" // Default Temporal server address
		}
	}

	logger.Info(nil, "Connecting to Temporal server", map[string]interface{}{
		"address": temporalAddr,
	})

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort:  temporalAddr,
		Namespace: "default",
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create Temporal client: %w", err)
	}

	// Create worker
	w := worker.New(c, TaskQueue, worker.Options{})

	// Register workflows
	w.RegisterWorkflow(workflows.ExampleWorkflow)
	w.RegisterWorkflow(workflows.CleanupWorkflow)

	// Register activities
	w.RegisterActivity(activities.ProcessMessageActivity)
	w.RegisterActivity(activities.SendNotificationActivity)
	w.RegisterActivity(activities.CleanupOldDataActivity)

	logger.Info(nil, "Temporal worker initialized", map[string]interface{}{
		"task_queue": TaskQueue,
	})

	return &WorkerService{
		client: c,
		worker: w,
	}, nil
}

// Start starts the Temporal worker and blocks until interrupted
func (ws *WorkerService) Start() error {
	logger.Info(nil, "Starting Temporal worker", nil)

	// Start worker in background
	err := ws.worker.Start()
	if err != nil {
		return fmt.Errorf("unable to start worker: %w", err)
	}

	logger.Info(nil, "Temporal worker started successfully", map[string]interface{}{
		"task_queue": TaskQueue,
		"status":     "running",
	})

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info(nil, "Received shutdown signal, stopping worker", nil)

	return nil
}

// Stop stops the Temporal worker
func (ws *WorkerService) Stop() {
	logger.Info(nil, "Stopping Temporal worker", nil)

	if ws.worker != nil {
		ws.worker.Stop()
	}

	if ws.client != nil {
		ws.client.Close()
	}

	logger.Info(nil, "Temporal worker stopped", nil)
}
