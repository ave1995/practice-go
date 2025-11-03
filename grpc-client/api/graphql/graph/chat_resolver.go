package graph

import (
	"context"
	"fmt"

	"github.com/ave1995/practice-go/grpc-client/connector/chat"
	"github.com/ave1995/practice-go/grpc-client/domain/model"
	"github.com/ave1995/practice-go/proto"
)

type ChatResolver struct {
	grpcConn *chat.Connector
}

func NewChatResolver(grpcConn *chat.Connector) *ChatResolver {
	return &ChatResolver{grpcConn: grpcConn}
}

func (r *ChatResolver) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	req := &proto.SendMessageRequest{
		Message: &proto.Message{
			Text: text,
		},
	}

	msg, err := r.grpcConn.SendMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("send message via gRPC: %w", err)
	}

	return model.FromGRPCMessage(msg), nil
}
