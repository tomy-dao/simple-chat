package workflows

import (
	"local/job/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// ExampleWorkflowInput defines the input for the example workflow
type ExampleWorkflowInput struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}

// ExampleWorkflowResult defines the result of the example workflow
type ExampleWorkflowResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ExampleWorkflow is a sample workflow that demonstrates Temporal usage
func ExampleWorkflow(ctx workflow.Context, input ExampleWorkflowInput) (*ExampleWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ExampleWorkflow started", "UserID", input.UserID, "Message", input.Message)

	// Set activity options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Execute activity
	var activityResult string
	err := workflow.ExecuteActivity(ctx, activities.ProcessMessageActivity, input.Message).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Activity failed", "Error", err)
		return &ExampleWorkflowResult{
			Success: false,
			Message: "Failed to process message: " + err.Error(),
		}, err
	}

	logger.Info("ExampleWorkflow completed", "Result", activityResult)

	return &ExampleWorkflowResult{
		Success: true,
		Message: activityResult,
	}, nil
}
