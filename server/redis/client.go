package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"rate-limiter/server/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewClient(host, port string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	// ping redis
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

func (r *Client) SetRateLimitRule(ctx context.Context, apiPath string, limit *models.Ratelimit) error {
	key := fmt.Sprintf("limit:%s", apiPath)
	return r.client.Set(ctx, key, limit, 0).Err()
}

func (r *Client) GetRateLimitRule(ctx context.Context, apiPath string) (*models.Ratelimit, error) {
	key := fmt.Sprintf("limit:%s", apiPath)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var limit models.Ratelimit
	if err := json.Unmarshal([]byte(val), &limit); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rate limit: %w", err)
	}
	return &limit, nil
}

func (r *Client) IncrementAndCheck(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	count := incr.Val()
	return count <= int64(limit), nil
}
