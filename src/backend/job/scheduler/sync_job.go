package scheduler

import (
	"fmt"
	"local/util/logger"
)

// RunSyncJob runs the sync job
func RunSyncJob() {
	logger.Info(nil, "Running sync job")
	// Sync job logic
	fmt.Println("Sync job executed")
}
