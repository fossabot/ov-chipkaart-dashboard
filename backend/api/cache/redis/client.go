package redis

import (
	"context"
	"time"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/cache"
	"github.com/go-redis/redis/v8"
)

// Client is a redis client
type Client struct {
	db *redis.Client
}

// Options are the params for initializing the class
type Options struct {
	Address  string
	Password string
	DB       int
}

// NewClient creates a new version of the cache
func NewClient(options Options) cache.Cache {
	return &Client{db: redis.NewClient(&redis.Options{
		Addr:     options.Address,
		Password: options.Password,
		DB:       options.DB,
	})}
}

// Set is responsible for setting a value in the cache
func (client *Client) Set(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return client.db.Set(ctx, key, value, expiration).Err()
}

// Get gets a value from the cache
func (client *Client) Get(key string) (result string, err error) {
	result, err = client.db.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return result, cache.ErrCacheMiss
	}

	if err != nil {
		return result, err
	}

	return result, err
}
