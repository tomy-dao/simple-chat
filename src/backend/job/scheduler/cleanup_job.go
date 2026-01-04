package scheduler

import (
	"fmt"
	"local/util/logger"
)

// RunCleanupJob runs the cleanup job
func RunCleanupJob() {
	logger.Info(nil, "Running cleanup job")
	// Cleanup job logic
	fmt.Println("Cleanup job executed")
}
