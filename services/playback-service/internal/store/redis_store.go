package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/freedom-music/playback-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(client *redis.Client, ttl time.Duration) *RedisStore {
	return &RedisStore{client: client, ttl: ttl}
}

func (s *RedisStore) Save(ctx context.Context, session *domain.PlaybackSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key(session.UserID), data, s.ttl).Err()
}

func (s *RedisStore) Get(ctx context.Context, userID int64) (*domain.PlaybackSession, error) {
	val, err := s.client.Get(ctx, key(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var sess domain.PlaybackSession
	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func key(userID int64) string {
	return fmt.Sprintf("playback:%d", userID)
}
