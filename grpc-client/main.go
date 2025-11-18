package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ave1995/practice-go/grpc-client/api/graphql"
	"github.com/ave1995/practice-go/grpc-client/api/graphql/graph"
	"github.com/ave1995/practice-go/grpc-client/config"
	"github.com/ave1995/practice-go/grpc-client/connector/chat"
	"github.com/ave1995/practice-go/grpc-client/service/message"
	"github.com/ave1995/practice-go/utils"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger := utils.NewInfoLogger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	chatConnector, err := chat.NewChatConnector(cfg.ChatConnectorConfig())
	if err != nil {
		logger.Error("grpc.NewConnector", utils.SlogError(err))
		os.Exit(1)
	}

	chatService := message.NewService(chatConnector)

	chatResolver := graph.NewResolver(chatService)

	graphql.RunGraphQLServer(ctx, cfg, logger, chatResolver)

	logger.Info("All services have stopped. Application exiting.")

	if errors.Is(ctx.Err(), context.Canceled) {
		fmt.Println("Graceful shutdown successful.")
	} else {
		fmt.Println("Application exited unexpectedly.")
	}
}
