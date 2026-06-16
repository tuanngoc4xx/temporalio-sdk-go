package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

// ActivityTaskHandler simulates the internal activity handler logic
type ActivityTaskHandler struct {
	heartbeatScheduler interface{ Flush() }
}

// CompleteActivity handles the completion of an activity, ensuring heartbeats are flushed
func (h *ActivityTaskHandler) CompleteActivity(ctx context.Context, result interface{}) error {
	// 1. Flush heartbeats synchronously before completion
	h.heartbeatScheduler.Flush()

	// 2. Check if context was already canceled due to timeout
	if errors.Is(ctx.Err(), context.Canceled) {
		return fmt.Errorf("activity already timed out: %w", ctx.Err())
	}

	// 3. Attempt to report completion to server
	err := reportCompletionToServer(result)
	if err != nil {
		// 4. Graceful error handling for race conditions
		if isTimeoutError(err) {
			log.Printf("Warning: activity completed but server already timed it out: %v", err)
			return nil // Discard result as server has already moved on
		}
		return err
	}

	return nil
}

func reportCompletionToServer(result interface{}) error {
	// Mock implementation of RPC call
	return nil
}

func isTimeoutError(err error) bool {
	return err.Error() == "ActivityTaskTimedOut"
}

func main() {
	fmt.Println("Temporal SDK Activity Handler logic updated.")
}