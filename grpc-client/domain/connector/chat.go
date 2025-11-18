package connector

import (
	"context"

	"github.com/ave1995/practice-go/grpc-client/domain/model"
)

type ChatConnector interface {
	SendMessage(ctx context.Context, text string) (*model.Message, error)
}
