package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/brianwu291/go-learn/cache"
)

type (
	Client struct {
		client *redis.Client
	}

	Pipeline struct {
		pipeline redis.Pipeliner
		cmds     []cache.PipelineCmd
	}

	pipelineIncrCmd struct {
		cmd *redis.IntCmd
	}

	pipelineTTLCmd struct {
		cmd *redis.DurationCmd
	}
)

func NewClient(cfg *cache.Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

func (c *Client) Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error) {
	return c.client.Eval(ctx, script, keys, args).Result()
}

func (c *Client) Pipeline() cache.Pipeline {
	return &Pipeline{
		pipeline: c.client.Pipeline(),
		cmds:     make([]cache.PipelineCmd, 0),
	}
}

func (p *Pipeline) Incr(ctx context.Context, key string) cache.PipelineCmd {
	cmd := &pipelineIncrCmd{cmd: p.pipeline.Incr(ctx, key)}
	p.cmds = append(p.cmds, cmd)
	return cmd
}

func (p *Pipeline) TTL(ctx context.Context, key string) cache.PipelineDurationCmd {
	return &pipelineTTLCmd{cmd: p.pipeline.TTL(ctx, key)}
}

func (p *Pipeline) Exec(ctx context.Context) error {
	_, err := p.pipeline.Exec(ctx)
	return err
}

func (c *pipelineIncrCmd) Val() int64 {
	return c.cmd.Val()
}

func (c *pipelineIncrCmd) Err() error {
	return c.cmd.Err()
}

func (c *pipelineTTLCmd) Val() time.Duration {
	return c.cmd.Val()
}

func (c *pipelineTTLCmd) Err() error {
	return c.cmd.Err()
}
