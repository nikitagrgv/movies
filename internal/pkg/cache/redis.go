package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.TODO()).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}

func GetOrSet[T any](ctx context.Context, client *redis.Client, key string, ttl time.Duration, fetch func() (T, error)) (T, error) {
	var result T

	raw, err := client.Get(ctx, key).Result()
	if err == nil {
		if jsonErr := json.Unmarshal([]byte(raw), &result); jsonErr == nil {
			return result, nil
		}
		// Error, fallback to fetch
	}

	result, err = fetch()
	if err != nil {
		return result, err
	}

	if bytes, marshalErr := json.Marshal(result); marshalErr == nil {
		client.Set(ctx, key, bytes, ttl)
	}

	return result, nil
}
