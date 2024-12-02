package repository

import (
	"context"
	"encoding/json"
	"time"

	"auth-service/internal/domain"

	"github.com/go-redis/redis/v8"
)

type SessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository(redis *redis.Client) *SessionRepository {
	return &SessionRepository{
		redis: redis,
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := "session:" + session.ID
	ttl := time.Until(session.ExpiresAt)

	return r.redis.Set(ctx, key, sessionData, ttl).Err()
}

func (r *SessionRepository) GetByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	key := "session:" + sessionID
	sessionData, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var session domain.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	key := "session:" + sessionID
	return r.redis.Del(ctx, key).Err()
}

func (r *SessionRepository) Validate(ctx context.Context, sessionID string) (*domain.Session, error) {
	key := "session:" + sessionID

	sessionData, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var session domain.Session
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = r.redis.Del(ctx, key).Err()
		return nil, nil
	}

	return &session, nil
}
