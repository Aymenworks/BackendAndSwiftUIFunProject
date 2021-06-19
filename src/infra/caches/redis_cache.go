package caches

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type client struct{}

func NewRedisCache() Cache {
	return &client{}
}

var (
	// TODO: use env variable
	redisClt = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)

func (c *client) Get(ctx context.Context, k string) (interface{}, error) {
	val, err := redisClt.Get(ctx, k).Result()
	switch {
	case err == redis.Nil:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return val, nil
	}
}

func (c *client) Set(ctx context.Context, k string, v interface{}, d time.Duration) error {
	return redisClt.Set(ctx, k, v, d).Err()
}

func (c *client) Delete(ctx context.Context, k string) error {
	return redisClt.Del(ctx, k).Err()
}

func (c *client) Ping() error {
	return redisClt.Ping(context.Background()).Err()
}
