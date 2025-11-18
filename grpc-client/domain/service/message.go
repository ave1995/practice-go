package service

import (
	"context"

	"github.com/ave1995/practice-go/grpc-client/domain/model"
)

type MessageService interface {
	Send(ctx context.Context, text string) (*model.Message, error)
}
