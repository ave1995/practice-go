package chat

import (
	"context"
	"fmt"

	"github.com/ave1995/practice-go/grpc-client/config"
	"github.com/ave1995/practice-go/grpc-client/domain/connector"
	"github.com/ave1995/practice-go/grpc-client/domain/model"
	"github.com/ave1995/practice-go/proto"
	ugrpc "github.com/ave1995/practice-go/utils/grpc"
)

var _ connector.ChatConnector = (*Connector)(nil)

type Connector struct {
	chatClient proto.ChatServiceClient
}

func NewChatConnector(config config.ChatClientConfig) (*Connector, error) {
	grpcConn, err := ugrpc.NewConnector(config.GRPCConfig)
	if err != nil {
		return nil, fmt.Errorf("create grpc connector: %w", err)
	}

	fmt.Printf("Attempting to connect to: %s\n", config.GRPCConfig.Address)

	if err := ugrpc.EnsureConnected(grpcConn, config.GRPCConfig.Timeout); err != nil {
		return nil, fmt.Errorf("gRPC server connectivity check failed: %w", err)
	}

	chatClient := proto.NewChatServiceClient(grpcConn)

	return &Connector{
		chatClient: chatClient,
	}, nil
}

func (c Connector) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	req := &proto.SendMessageRequest{
		Message: &proto.Message{
			Text: text,
		},
	}

	msg, err := c.chatClient.SendMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &model.Message{
		ID:   msg.Id,
		Text: msg.Message,
	}, nil

}
