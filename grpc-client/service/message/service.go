package message

import (
	"context"
	"fmt"

	"github.com/ave1995/practice-go/grpc-client/domain/connector"
	"github.com/ave1995/practice-go/grpc-client/domain/model"
	"github.com/ave1995/practice-go/grpc-client/domain/service"
)

var _ service.MessageService = (*Service)(nil)

type Service struct {
	chatConnector connector.ChatConnector
}

func NewService(chatConnector connector.ChatConnector) *Service {
	return &Service{chatConnector: chatConnector}
}

func (s Service) Send(ctx context.Context, text string) (*model.Message, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}
	return s.chatConnector.SendMessage(ctx, text)
}
