package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"rate-limiter/server/models"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewClient(host, port string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	return &Client{client: client}, nil
}

func (r *Client) SetRateLimit(ctx context.Context, apiPath string, limit *models.Ratelimit) error {

}

func (r *Client) GetRateLimit(ctx context.Context, apiPath string) (*models.Ratelimit, error) {
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
