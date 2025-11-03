package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ave1995/practice-go/grpc-client/api/graphql/graph"
	"github.com/ave1995/practice-go/grpc-client/config"
	"github.com/ave1995/practice-go/grpc-client/connector/chat"
	"github.com/ave1995/practice-go/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

const DefaultPort = "8080"

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := utils.NewInfoLogger()

	chatConnector, err := chat.NewChatConnector(cfg.ChatConnectorConfig())
	if err != nil {
		logger.Error("grpc.NewConnector", utils.SlogError(err))
		os.Exit(1)
	}

	chatResolver := graph.NewChatResolver(chatConnector)

	mux := http.NewServeMux()
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: chatResolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	server := &http.Server{
		Addr:    ":" + DefaultPort,
		Handler: mux,
	}

	go func() {
		logger.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", defaultPort))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http.ListenAndServe", utils.SlogError(err))
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Info("Received signal. Shutting down gracefully...", "signal", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server.Shutdown", utils.SlogError(err))
	}

	cancel()
	logger.Info("Server stopped cleanly.")
}
