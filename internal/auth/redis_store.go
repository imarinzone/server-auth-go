package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore wraps another Store and adds caching
type RedisStore struct {
	client   *redis.Client
	next     Store
	cacheTTL time.Duration
}

// NewRedisStore creates a new Redis store
func NewRedisStore(addr, password string, db int, next Store) (*RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisStore{
		client:   rdb,
		next:     next,
		cacheTTL: 10 * time.Minute,
	}, nil
}

// VerifyCredentials checks Redis first, then falls back to the underlying store
func (s *RedisStore) VerifyCredentials(clientID, clientSecret string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("client:%s", clientID)

	// Check Cache
	val, err := s.client.Get(ctx, key).Result()
	if err == nil {
		// Cache hit
		return val == clientSecret, nil
	} else if err != redis.Nil {
		// Redis error, log it but continue to DB
		fmt.Printf("Redis error: %v\n", err)
	}

	// Cache Miss - Check DB
	valid, err := s.next.VerifyCredentials(clientID, clientSecret)
	if err != nil {
		return false, err
	}

	if valid {
		// Update Cache
		if err := s.client.Set(ctx, key, clientSecret, s.cacheTTL).Err(); err != nil {
			fmt.Printf("Failed to update cache: %v\n", err)
		}
	}

	return valid, nil
}
