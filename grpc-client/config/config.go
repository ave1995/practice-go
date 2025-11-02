package config

import (
	"fmt"
	"time"

	"github.com/ave1995/practice-go/utils/grpc"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ChatClientAddress        string        `envconfig:"CHAT_CLIENT_ADDRESS"`
	ChatClientUseTLS         bool          `envconfig:"CHAT_CLIENT_USE_TLS"`
	ChatClientCertFile       string        `envconfig:"CHAT_CLIENT_CERT_FILE"`
	ChatClientTimeout        time.Duration `envconfig:"CHAT_CLIENT_TIMEOUT"`
	ChatClientMaxRetries     int           `envconfig:"CHAT_CLIENT_MAX_RETRIES"`
	ChatClientKeepAlive      bool          `envconfig:"CHAT_CLIENT_KEEPALIVE"`
	ChatClientEnableRetry    bool          `envconfig:"CHAT_CLIENT_ENABLE_RETRY"`
	ChatClientMaxMessageSize int           `envconfig:"CHAT_CLIENT_MAX_MESSAGE_SIZE"`
}

const EnvPrefix = ""

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process(EnvPrefix, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to process env vars: %w", err)
	}

	return cfg, nil
}

func (c Config) ChatConnectorConfig() ChatClientConfig {
	return ChatClientConfig{
		grpc.Config{
			Address:        c.ChatClientAddress,
			UseTLS:         c.ChatClientUseTLS,
			CertFile:       c.ChatClientCertFile,
			Timeout:        c.ChatClientTimeout,
			MaxRetries:     c.ChatClientMaxRetries,
			KeepAlive:      c.ChatClientKeepAlive,
			EnableRetry:    c.ChatClientEnableRetry,
			MaxMessageSize: c.ChatClientMaxMessageSize,
		},
	}
}
