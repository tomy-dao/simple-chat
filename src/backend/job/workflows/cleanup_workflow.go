package workflows

import (
	"local/job/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CleanupWorkflowInput defines the input for the cleanup workflow
type CleanupWorkflowInput struct {
	OlderThanDays int `json:"older_than_days"`
}

// CleanupWorkflowResult defines the result of the cleanup workflow
type CleanupWorkflowResult struct {
	Success      bool `json:"success"`
	DeletedCount int  `json:"deleted_count"`
}

// CleanupWorkflow is a workflow that cleans up old data
func CleanupWorkflow(ctx workflow.Context, input CleanupWorkflowInput) (*CleanupWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("CleanupWorkflow started", "OlderThanDays", input.OlderThanDays)

	// Set activity options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Execute cleanup activity
	var deletedCount int
	err := workflow.ExecuteActivity(ctx, activities.CleanupOldDataActivity, input.OlderThanDays).Get(ctx, &deletedCount)
	if err != nil {
		logger.Error("Cleanup activity failed", "Error", err)
		return &CleanupWorkflowResult{
			Success:      false,
			DeletedCount: 0,
		}, err
	}

	logger.Info("CleanupWorkflow completed", "DeletedCount", deletedCount)

	return &CleanupWorkflowResult{
		Success:      true,
		DeletedCount: deletedCount,
	}, nil
}
