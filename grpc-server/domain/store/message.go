package store

import (
	"context"

	"github.com/ave1995/practise-go/grpc-server/domain/model"
)

type MessageStore interface {
	Fetch(ctx context.Context, id model.MessageID) (*model.Message, error)
	Create(ctx context.Context, text string) (*model.Message, error)
}
