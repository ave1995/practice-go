package model

import "github.com/ave1995/practice-go/proto"

type Message struct {
	ID   string
	Text string
}

func FromGRPCMessage(msg *proto.SendMessageResponse) *Message {
	if msg == nil {
		return nil
	}
	return &Message{
		ID:   msg.Id,
		Text: msg.Message,
	}
}
