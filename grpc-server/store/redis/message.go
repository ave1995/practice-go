package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ave1995/practice-go/grpc-server/domain/model"
	"github.com/ave1995/practice-go/grpc-server/domain/store"
	"github.com/ave1995/practice-go/utils"
	"github.com/redis/go-redis/v9"
)

var _ store.MessageStore = (*MessageStore)(nil)

type MessageStore struct {
	logger       *slog.Logger
	messageStore store.MessageStore
	redisClient  *redis.Client
}

func NewMessageStore(logger *slog.Logger, store store.MessageStore, addr string, password string, db int) *MessageStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &MessageStore{
		logger:       logger,
		messageStore: store,
		redisClient:  rdb,
	}
}

func (m *MessageStore) Fetch(ctx context.Context, id model.MessageID) (*model.Message, error) {
	key := id.String()

	data, err := m.redisClient.Get(ctx, key).Bytes()
	if err == nil {
		var msg model.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached message: %w", err)
		}
		return &msg, nil
	}

	if !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get message from redis: %w", err)
	}

	msg, err := m.messageStore.Fetch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message from store: %w", err)
	}

	data, err = json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := m.redisClient.Set(ctx, key, data, 0).Err(); err != nil {
		m.logger.Warn("failed to set redis cache", utils.SlogError(err))
	}

	return msg, nil
}

func (m *MessageStore) Create(ctx context.Context, text string) (*model.Message, error) {
	msg, err := m.messageStore.Create(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	key := msg.ID.String()

	if err := m.redisClient.Set(ctx, key, data, 0).Err(); err != nil {
		m.logger.Warn("failed to set redis cache", utils.SlogError(err))
	}

	return msg, nil
}
