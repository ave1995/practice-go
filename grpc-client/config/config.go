package config

import (
	"fmt"
	"time"

	"github.com/ave1995/practice-go/utils/grpc"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ChatServerAddress        string        `envconfig:"CHAT_SERVER_ADDRESS"`
	ChatServerUseTLS         bool          `envconfig:"CHAT_SERVER_USE_TLS"`
	ChatServerCertFile       string        `envconfig:"CHAT_SERVER_CERT_FILE"`
	ChatServerTimeout        time.Duration `envconfig:"CHAT_SERVER_TIMEOUT"`
	ChatServerMaxRetries     int           `envconfig:"CHAT_SERVER_MAX_RETRIES"`
	ChatServerKeepAlive      bool          `envconfig:"CHAT_SERVER_KEEPALIVE"`
	ChatServerEnableRetry    bool          `envconfig:"CHAT_SERVER_ENABLE_RETRY"`
	ChatServerMaxMessageSize int           `envconfig:"CHAT_SERVER_MAX_MESSAGE_SIZE"`
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
			Address:        c.ChatServerAddress,
			UseTLS:         c.ChatServerUseTLS,
			CertFile:       c.ChatServerCertFile,
			Timeout:        c.ChatServerTimeout,
			MaxRetries:     c.ChatServerMaxRetries,
			KeepAlive:      c.ChatServerKeepAlive,
			EnableRetry:    c.ChatServerEnableRetry,
			MaxMessageSize: c.ChatServerMaxMessageSize,
		},
	}
}
