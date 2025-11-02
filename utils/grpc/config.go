package grpc

import "time"

type Config struct {
	Address        string
	UseTLS         bool
	CertFile       string
	Timeout        time.Duration
	MaxRetries     int
	KeepAlive      bool
	EnableRetry    bool
	MaxMessageSize int
}
