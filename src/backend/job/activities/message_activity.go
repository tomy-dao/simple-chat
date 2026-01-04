package activities

import (
	"context"
	"fmt"
	"strings"

	"go.temporal.io/sdk/activity"
)

// ProcessMessageActivity processes a message and returns the result
func ProcessMessageActivity(ctx context.Context, message string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ProcessMessageActivity started", "Message", message)

	// Simulate some processing
	processedMessage := fmt.Sprintf("Processed: %s (length: %d, uppercase: %s)",
		message,
		len(message),
		strings.ToUpper(message),
	)

	logger.Info("ProcessMessageActivity completed", "Result", processedMessage)
	return processedMessage, nil
}

// SendNotificationActivity sends a notification
func SendNotificationActivity(ctx context.Context, userID uint, message string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("SendNotificationActivity started", "UserID", userID, "Message", message)

	// TODO: Implement actual notification sending logic
	// For now, just log it
	logger.Info("Notification sent successfully", "UserID", userID)

	return nil
}

// CleanupOldDataActivity cleans up old data from the database
func CleanupOldDataActivity(ctx context.Context, olderThanDays int) (int, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CleanupOldDataActivity started", "OlderThanDays", olderThanDays)

	// TODO: Implement actual cleanup logic
	// For now, return a mock count
	deletedCount := 42

	logger.Info("CleanupOldDataActivity completed", "DeletedCount", deletedCount)
	return deletedCount, nil
}
