package graph

import (
	"github.com/ave1995/practice-go/grpc-client/domain/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	messageService service.MessageService
}

func NewResolver(messageService service.MessageService) *Resolver {
	return &Resolver{messageService: messageService}
}
