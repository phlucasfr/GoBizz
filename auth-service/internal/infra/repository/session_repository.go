package repository

import (
	"context"
	"encoding/json"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type SessionRepository struct {
	redis *redis.Client
}

// NewSessionRepository creates a new instance of SessionRepository with the provided Redis client.
// It initializes the repository to interact with the Redis database for session management.
//
// Parameters:
//   - redis: A pointer to a redis.Client instance used to interact with the Redis database.
//
// Returns:
//   - A pointer to a newly created SessionRepository instance.
func NewSessionRepository(redis *redis.Client) *SessionRepository {
	return &SessionRepository{
		redis: redis,
	}
}

// Create stores a session in the Redis database with a specified time-to-live (TTL).
// It serializes the session object into JSON format and saves it using the session ID as the key.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellation signals.
//   - session: A pointer to the Session object containing session details to be stored.
//
// Returns:
//   - error: An error if the session could not be serialized or stored in Redis, otherwise nil.
func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		logger.Log.Error("Failed to marshal session data", zap.Error(err))
		return err
	}

	key := "session:" + session.ID
	ttl := time.Until(session.ExpiresAt)

	logger.Log.Info("Storing session in Redis", zap.String("key", key), zap.Duration("ttl", ttl))
	return r.redis.Set(ctx, key, sessionData, ttl).Err()
}

// GetByID retrieves a session by its ID from the Redis datastore.
// It takes a context and a session ID as input parameters and returns a pointer
// to a domain.Session object and an error. If the session is not found, it returns
// nil for both the session and the error. If an error occurs during the retrieval
// or unmarshalling process, it returns the error.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellation signals.
//   - sessionID: The unique identifier of the session to retrieve.
//
// Returns:
//   - *domain.Session: A pointer to the session object if found, or nil if not found.
//   - error: An error if any issue occurs during the retrieval or unmarshalling process.
func (r *SessionRepository) GetByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	key := "session:" + sessionID
	sessionData, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Log.Info("Session not found", zap.String("key", key))
		return nil, nil
	} else if err != nil {
		logger.Log.Error("Failed to get session from Redis", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	var session domain.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		logger.Log.Error("Failed to unmarshal session data", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Session retrieved from Redis", zap.String("key", key), zap.String("sessionID", session.ID))
	return &session, nil
}

// Delete removes a session from the Redis datastore based on the provided session ID.
// It constructs the Redis key using the session ID and deletes the corresponding entry.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellation signals.
//   - sessionID: The unique identifier of the session to be deleted.
//
// Returns:
//   - error: An error if the deletion operation fails, otherwise nil.
func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	key := "session:" + sessionID
	logger.Log.Info("Deleting session from Redis", zap.String("key", key))
	return r.redis.Del(ctx, key).Err()
}

// Validate checks the validity of a session based on its session ID.
// It retrieves the session data from Redis using the provided session ID,
// unmarshals it into a Session object, and verifies if the session has expired.
// If the session is expired, it deletes the session from Redis and returns nil.
// If the session is valid, it returns the Session object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellations.
//   - sessionID: The unique identifier of the session to validate.
//
// Returns:
//   - *domain.Session: The session object if it is valid, or nil if it is invalid or not found.
//   - error: An error if there is an issue retrieving or unmarshaling the session data.
func (r *SessionRepository) Validate(ctx context.Context, sessionID string) (*domain.Session, error) {
	key := "session:" + sessionID

	sessionData, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Log.Info("Session not found", zap.String("key", key))
		return nil, nil
	} else if err != nil {
		logger.Log.Error("Failed to get session from Redis", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	var session domain.Session
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		logger.Log.Error("Failed to unmarshal session data", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = r.redis.Del(ctx, key).Err()
		logger.Log.Info("Session expired and deleted", zap.String("key", key))
		return nil, nil
	}

	logger.Log.Info("Session is valid", zap.String("key", key), zap.String("sessionID", session.ID))
	return &session, nil
}
