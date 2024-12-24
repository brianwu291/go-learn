package cache

import (
	"context"
	"errors"
	"time"
)

// Common cache errors
var (
	// ErrKeyNotFound is returned when a key doesn't exist
	ErrKeyNotFound = errors.New("key not found in cache")
	// ErrKeyExpired is returned when a key has expired
	ErrKeyExpired = errors.New("key has expired")
	// ErrConnectionFailed is returned when cache connection fails
	ErrConnectionFailed = errors.New("failed to connect to cache")
	// ErrInvalidValue is returned when value is invalid or corrupted
	ErrInvalidValue = errors.New("invalid cache value")
)

type (
	CacheError interface {
		error
		IsCacheError() bool
	}

	KeyNotFoundError struct {
		Key string
	}

	ConnectionError struct {
		Err error
	}

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

func (e *KeyNotFoundError) Error() string {
	if e.Key == "" {
		return ErrKeyNotFound.Error()
	}
	return "key not found in cache: " + e.Key
}

func (e *KeyNotFoundError) IsCacheError() bool {
	return true
}

func (e *KeyNotFoundError) Is(target error) bool {
	return target == ErrKeyNotFound
}

func (e *ConnectionError) Error() string {
	if e.Err == nil {
		return ErrConnectionFailed.Error()
	}
	return "cache connection failed: " + e.Err.Error()
}

func (e *ConnectionError) IsCacheError() bool {
	return true
}

func (e *ConnectionError) Is(target error) bool {
	return target == ErrConnectionFailed
}

// Helper functions to check error types
func IsKeyNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrKeyNotFound)
}

func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrConnectionFailed)
}

// Helper functions to create errors
func NewKeyNotFoundError(key string) error {
	return &KeyNotFoundError{Key: key}
}

func NewConnectionError(err error) error {
	return &ConnectionError{Err: err}
}
