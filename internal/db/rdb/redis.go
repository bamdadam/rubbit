package rdb

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDb struct {
	Client *redis.Client
}

func New(ctx context.Context) (*RedisDb, error) {
	const t = 10
	timeout, cancel := context.WithTimeout(ctx, t*time.Second)
	defer cancel()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := client.Ping(timeout).Err(); err != nil {
		return nil, err
	}
	return &RedisDb{
		Client: client,
	}, nil
}
