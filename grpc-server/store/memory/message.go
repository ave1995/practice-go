package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ave1995/practice-go/grpc-server/domain/model"
	"github.com/ave1995/practice-go/grpc-server/domain/store"
)

var _ store.MessageStore = (*MessageStore)(nil)

type MessageStore struct {
	messageStore store.MessageStore
	mu           sync.RWMutex
	items        map[model.MessageID]*model.Message
}

func NewMessageStore(store store.MessageStore) *MessageStore {
	return &MessageStore{
		messageStore: store,
		items:        make(map[model.MessageID]*model.Message),
	}
}

func (m *MessageStore) Fetch(ctx context.Context, id model.MessageID) (*model.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	item, found := m.items[id]
	m.mu.RUnlock()

	if found {
		return item, nil
	}

	msg, err := m.messageStore.Fetch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message: %w", err)
	}

	m.mu.Lock()
	m.items[id] = msg
	m.mu.Unlock()

	return msg, nil
}

func (m *MessageStore) Create(ctx context.Context, text string) (*model.Message, error) {
	msg, err := m.messageStore.Create(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	m.mu.Lock()
	m.items[msg.ID] = msg
	m.mu.Unlock()

	return msg, nil
}
