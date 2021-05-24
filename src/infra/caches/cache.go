package caches

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, k string) (interface{}, error)
	Set(ctx context.Context, k string, v interface{}, d time.Duration) error
	Ping() error
}
