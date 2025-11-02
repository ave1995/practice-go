package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

func NewConnector(config Config) (*grpc.ClientConn, error) {
	opts, err := buildDialOptions(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build dial options: %w", err)
	}

	return grpc.NewClient(config.Address, opts...)
}

// buildDialOptions constructs gRPC dial options based on configuration
func buildDialOptions(config Config) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	// TLS configuration
	if config.UseTLS {
		if config.CertFile != "" {
			creds, err := credentials.NewClientTLSFromFile(config.CertFile, "")
			if err != nil {
				return nil, fmt.Errorf("failed to create tls client for %s: %w", config.Address, err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			// Use system cert pool
			creds := credentials.NewTLS(nil)
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
	} else {
		// Only for development/testing
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Keepalive configuration
	if config.KeepAlive {
		kaParams := keepalive.ClientParameters{
			Time:                10 * time.Second, // Send keepalive every 10 seconds
			Timeout:             3 * time.Second,  // Wait 3 seconds for response
			PermitWithoutStream: true,             // Send pings even without active streams
		}
		opts = append(opts, grpc.WithKeepaliveParams(kaParams))
	}

	// Max message size
	opts = append(opts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(config.MaxMessageSize),
			grpc.MaxCallSendMsgSize(config.MaxMessageSize),
		),
	)

	// Retry configuration
	if config.EnableRetry {
		opts = append(opts, grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{"service": ""}],
				"retryPolicy": {
					"maxAttempts": 3,
					"initialBackoff": "0.1s",
					"maxBackoff": "1s",
					"backoffMultiplier": 2,
					"retryableStatusCodes": ["UNAVAILABLE", "RESOURCE_EXHAUSTED"]
				}
			}]
		}`))
	}

	// Add interceptors
	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			loggingInterceptor(),
			retryInterceptor(config.MaxRetries),
		),
	)

	return opts, nil
}

// loggingInterceptor logs RPC calls
func loggingInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		// Log the call
		if err != nil {
			fmt.Printf("[gRPC] %s failed in %v: %v\n", method, duration, err)
		} else {
			fmt.Printf("[gRPC] %s succeeded in %v\n", method, duration)
		}

		return err
	}
}

// retryInterceptor implements retry logic with exponential backoff
func retryInterceptor(maxRetries int) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		var lastErr error
		backoff := 100 * time.Millisecond

		for attempt := 0; attempt < maxRetries; attempt++ {
			if attempt > 0 {
				// Wait before retry
				select {
				case <-time.After(backoff):
					backoff *= 2 // Exponential backoff
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			lastErr = invoker(ctx, method, req, reply, cc, opts...)

			if lastErr == nil {
				return nil
			}

			// Check if error is retryable
			if !isRetryable(lastErr) {
				return lastErr
			}

			fmt.Printf("[gRPC] Retry attempt %d/%d for %s\n", attempt+1, maxRetries, method)
		}

		return lastErr
	}
}

func isRetryable(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Unavailable, codes.ResourceExhausted, codes.Aborted:
		return true
	default:
		return false
	}
}
