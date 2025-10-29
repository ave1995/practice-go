package service

import (
	"context"

	"github.com/ave1995/practice-go/grpc-server/domain/model"
)

type MessageService interface {
	Send(ctx context.Context, text string) (*model.Message, error)
	Fetch(ctx context.Context, id model.MessageID) (*model.Message, error)
	NewSubscriberWithCleanup() (*model.MessageSubscriber, func())
}
