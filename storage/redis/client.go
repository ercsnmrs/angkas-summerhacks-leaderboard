package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// Config represents the redis configuration
type Config struct {
	Host     string
	Password string
	Port     string
}

// Client represents the client for the internal redis package
type Client struct {
	*redis.Client
	config *Config
	logger *slog.Logger
}

// NewClient create new instance of redis client.
func New(config *Config, logger *slog.Logger) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),
	})

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		panic(err)
	}

	c := &Client{Client: redisClient, config: config, logger: logger}
	if _, err := c.Ping(); err != nil {
		return nil, err
	}
	return c, nil
}

// Ping checks redis connection
func (c *Client) Ping() (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = c.Client.Ping(ctx).Err(); err != nil {
		return false, fmt.Errorf("redis: could not ping client: %s", err)
	}
	return true, nil
}



// Close closes all connection.
func (c *Client) Close() error {
	return c.Client.Close()
}
