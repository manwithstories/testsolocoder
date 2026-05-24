package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"furniture-platform/pkg/config"
)

// Client Redis 客户端封装
type Client struct {
	rdb *redis.Client
}

// NewClient 创建 Redis 客户端
func NewClient(cfg *config.RedisConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}
	return &Client{rdb: rdb}, nil
}

// RDB 返回原始 Redis 客户端
func (c *Client) RDB() *redis.Client {
	return c.rdb
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.rdb.Close()
}
