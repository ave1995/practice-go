package chat

import (
	"context"
	"fmt"

	"github.com/ave1995/practice-go/grpc-client/config"
	"github.com/ave1995/practice-go/proto"
	ugrpc "github.com/ave1995/practice-go/utils/grpc"
	"google.golang.org/grpc"
)

var _ proto.ChatServiceClient = (*Connector)(nil)

type Connector struct {
	chatClient proto.ChatServiceClient
}

func NewChatConnector(config config.ChatClientConfig) (*Connector, error) {
	grpcConn, err := ugrpc.NewConnector(config.GRPCConfig)
	if err != nil {
		return nil, fmt.Errorf("create grpc connector: %w", err)
	}

	//
	//fmt.Printf("Attempting to connect to: %s\n", config.GRPCConfig.Address)
	//
	//if err := ugrpc.EnsureConnected(grpcConn, config.GRPCConfig.Timeout); err != nil {
	//	return nil, fmt.Errorf("gRPC server connectivity check failed: %w", err)
	//}

	chatClient := proto.NewChatServiceClient(grpcConn)

	return &Connector{
		chatClient: chatClient,
	}, nil
}

func (c Connector) SendMessage(ctx context.Context, in *proto.SendMessageRequest, opts ...grpc.CallOption) (*proto.SendMessageResponse, error) {
	return c.chatClient.SendMessage(ctx, in, opts...)
}

func (c Connector) GetMessage(
	ctx context.Context,
	in *proto.GetMessageRequest,
	opts ...grpc.CallOption,
) (*proto.GetMessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c Connector) Reader(ctx context.Context, in *proto.ReaderRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[proto.Message], error) {
	//TODO implement me
	panic("implement me")
}
