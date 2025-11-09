package graph

import "github.com/ave1995/practice-go/grpc-client/connector/chat"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	grpcConn *chat.Connector
}

func NewResolver(grpcConn *chat.Connector) *Resolver {
	return &Resolver{grpcConn: grpcConn}
}
