package model

import (
	"time"

	"github.com/ave1995/practice-go/proto"
)

type Message struct {
	ID        MessageID
	Text      string
	Timestamp time.Time
}

func (m Message) ToProto() *proto.Message {
	return &proto.Message{
		Text: m.Text,
	}
}
