package gormdb

import (
	"time"

	"github.com/ave1995/practise-go/grpc-server/domain/model"
	"github.com/google/uuid"
)

type message struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Text      string    `gorm:"type:text;not null"`
	Timestamp time.Time `gorm:"not null;autoCreateTime"`
}

func (m *message) ToDomain() *model.Message {
	return &model.Message{
		ID:        model.MessageID(m.ID),
		Text:      m.Text,
		Timestamp: m.Timestamp,
	}
}
