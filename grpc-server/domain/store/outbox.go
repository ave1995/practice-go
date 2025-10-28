package store

import (
	"context"

	"github.com/ave1995/practise-go/grpc-server/domain/model"
)

type OutboxStore interface {
	GetPendingEvents(ctx context.Context, eType model.EventType, limit int) ([]*model.OutboxEvent, error)
	MarkProcessed(ctx context.Context, id model.OutboxEventID) error
	MarkFailed(ctx context.Context, id model.OutboxEventID) error
}
