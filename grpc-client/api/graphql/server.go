package graphql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ave1995/practice-go/grpc-client/api/graphql/graph"
	"github.com/ave1995/practice-go/grpc-client/config"
	"github.com/ave1995/practice-go/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

func RunGraphQLServer(ctx context.Context, cfg config.Config, logger *slog.Logger, chatResolver *graph.Resolver) {
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
		Addr:    ":" + cfg.ServicePort,
		Handler: mux,
	}

	go func() {
		logger.Info(fmt.Sprintf("connect to http://localhost:%s/ for GraphQL playground", cfg.ServicePort))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http.ListenAndServe", utils.SlogError(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	logger.Info("Context cancelled. Starting HTTP server shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server forced shutdown after timeout", utils.SlogError(err))
	} else {
		logger.Info("HTTP server shut down successfully.")
	}
}
