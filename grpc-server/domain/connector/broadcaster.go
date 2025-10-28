package connector

import (
	"github.com/ave1995/practise-go/grpc-server/domain/model"
)

type Broadcaster interface {
	Broadcast(msg *model.Message) error
}
