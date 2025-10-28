package grpc

import (
	"context"
	"log/slog"

	"github.com/ave1995/practise-go/grpc-server/domain/model"
	"github.com/ave1995/practise-go/grpc-server/domain/service"
	"github.com/ave1995/practise-go/grpc-server/utils"
	"github.com/ave1995/practise-go/proto"
	"github.com/google/uuid"
)

type ChatServer struct {
	proto.ChatServiceServer
	logger         *slog.Logger
	messageService service.MessageService
}

func NewChatServer(logger *slog.Logger, messageService service.MessageService) *ChatServer {
	return &ChatServer{
		logger:         logger,
		messageService: messageService}
}

func (s *ChatServer) SendMessage(ctx context.Context, msg *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	created, err := s.messageService.Send(ctx, msg.Message.Text)
	if err != nil {
		return nil, err
	}
	s.logger.Info("message sent", "message", created.Text)

	return &proto.SendMessageResponse{Message: "Message stored successfully", Id: created.ID.String()}, nil
}

func (s *ChatServer) GetMessage(ctx context.Context, req *proto.GetMessageRequest) (*proto.GetMessageResponse, error) {
	u, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	id := model.MessageID(u)

	found, err := s.messageService.Fetch(ctx, id)
	if err != nil {
		return nil, err
	}

	return &proto.GetMessageResponse{
		Message: &proto.Message{
			Text: found.Text,
		},
	}, nil
}

func (s *ChatServer) Reader(_ *proto.ReaderRequest, srv proto.ChatService_ReaderServer) error {
	s.logger.Info("server stream opened")

	subscriber, cleanup := s.messageService.NewSubscriberWithCleanup()
	defer cleanup()

	for {
		select {
		case <-srv.Context().Done():
			s.logger.Info("server stream closed by disconnection of client: %v", "subscriber", subscriber)
			return nil

		case msg, open := <-subscriber.Messages():
			if !open {
				return nil
			}
			err := srv.Send(msg.ToProto())
			if err != nil {
				s.logger.Error("send error", utils.SlogError(err))
				return err
			}
		}
	}
}
