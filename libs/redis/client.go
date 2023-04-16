package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	c *goredis.Client
}

func NewClient(addr string) *Client {
	opts := &goredis.Options{
		Addr: addr,
	}
	client := goredis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return &Client{c: client}
}

func (c *Client) Set(ctx context.Context, key string, v any, exp time.Duration) *goredis.StatusCmd {
	return c.c.Set(ctx, key, v, exp)
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, bool, error) {
	v, err := c.c.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("client: Get: %w", err)
	}
	return v, true, nil
}
