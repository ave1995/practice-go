package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// EnsureConnected forces the connection to initiate if it is IDLE and then
// waits for the connection to reach the READY state within the specified timeout.
// This function implements the non-deprecated pattern for checking connection status
// and provides clean error handling for timeouts.
func EnsureConnected(grpcConn *grpc.ClientConn, timeout time.Duration) error {
	connCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	currentState := grpcConn.GetState()

	// 1. Force connection attempt if currently Idle (lazy connect strategy).
	// This kicks the gRPC channel out of its initial IDLE state.
	if currentState == connectivity.Idle {
		grpcConn.Connect()
		currentState = grpcConn.GetState() // State should now be CONNECTING
	}

	// 2. Loop until READY or timeout.
	for currentState != connectivity.Ready {

		// WaitForStateChange blocks until the state changes or the context is done.
		if !grpcConn.WaitForStateChange(connCtx, currentState) {

			// Connection failed/timed out or the context was cancelled.
			// The state is likely TRANSIENT_FAILURE. Close the channel to clean up resources.
			closeErr := grpcConn.Close()

			// Build the main error message reporting the failure.
			mainErr := fmt.Errorf(
				"gRPC connection failed to become READY within %s. Final state: %s",
				timeout,
				grpcConn.GetState(),
			)

			// If closing also failed, wrap both errors.
			if closeErr != nil {
				return fmt.Errorf("%w; failed to close connection after timeout: %v", mainErr, closeErr)
			}

			// Return the detailed connection error.
			return mainErr
		}

		// State changed. Get the new state and continue the loop check.
		currentState = grpcConn.GetState()
	}

	// If the loop finished, the connection is READY.
	return nil
}
