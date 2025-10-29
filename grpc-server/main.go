package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	lgprc "github.com/ave1995/practice-go/grpc-server/api/grpc"
	"github.com/ave1995/practice-go/grpc-server/service/message"
	"github.com/ave1995/practice-go/grpc-server/store/gormdb"
	"github.com/ave1995/practice-go/proto"
	"google.golang.org/grpc"

	"github.com/ave1995/practice-go/grpc-server/config"
	"github.com/ave1995/practice-go/grpc-server/factory"
	"github.com/ave1995/practice-go/grpc-server/utils"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fact := factory.NewFactory(ctx, cfg)
	logger := fact.Logger()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("net.listen", utils.SlogError(err))
		os.Exit(1)
	}

	messageService := fact.MessageService()
	messageService.Broadcast(ctx)

	grpcServer := grpc.NewServer()
	proto.RegisterChatServiceServer(grpcServer, lgprc.NewChatServer(logger, fact.MessageService()))

	outboxStore := gormdb.NewOutboxStore(fact.Database())

	processor := message.NewProcessor(logger, cfg.MessageProcessorConfig(), outboxStore, fact.KafkaProducer())
	processor.Start(ctx)

	go func() {
		logger.Info("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("grpcServer.serve", utils.SlogError(err))
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Info("Received signal. Shutting down gracefully...", "signal", sig)

	grpcServer.GracefulStop()

	fact.Close()

	cancel()

	logger.Info("Server stopped cleanly.")
}
