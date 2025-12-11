package cache

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Host     string
	Port     string
	Password string
}

type Client struct {
	client *redis.Client
}

func LoadConfig() Config {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	return Config{
		Addr:     addr,
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
}

func Connect(cfg Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
		Protocol: 2,
	})

	return &Client{rdb}
}

func (c *Client) RDB() *redis.Client {
	return c.client
}

func (c *Client) Close() error {
	return c.client.Close()
}
