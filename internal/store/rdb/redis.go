package rdb

import "github.com/redis/go-redis/v9"

type RedisStore struct {
	DB *redis.Client
}

func New(db *redis.Client) *RedisStore {
	return &RedisStore{
		DB: db,
	}
}
