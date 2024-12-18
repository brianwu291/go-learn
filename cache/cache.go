package cache

import (
	"context"
	"time"
)

type (
	PipelineCmd interface {
		Val() int64
		Err() error
	}

	PipelineDurationCmd interface {
		Val() time.Duration
		Err() error
	}

	Pipeline interface {
		Incr(ctx context.Context, key string) PipelineCmd
		TTL(ctx context.Context, key string) PipelineDurationCmd
		Exec(ctx context.Context) error
	}

	Client interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Incr(ctx context.Context, key string) (int64, error)
		TTL(ctx context.Context, key string) (time.Duration, error)
		Expire(ctx context.Context, key string, expiration time.Duration) error
		Pipeline() Pipeline
		Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error)
	}

	Config struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
)
