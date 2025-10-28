package connector

import (
	"context"

	"github.com/ave1995/practise-go/grpc-server/domain/model"
)

// Consumer defines the interface for reading messages from Kafka topics.
type Consumer interface {
	Read(ctx context.Context) (<-chan model.Message, error)
}
