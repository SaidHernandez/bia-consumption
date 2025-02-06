package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, bool, error)
	Set(ctx context.Context, key string, val interface{}, expireAfter time.Duration) error
	Clear(ctx context.Context, key string) (bool, error)
}
