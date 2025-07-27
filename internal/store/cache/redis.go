package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(addr string, pw string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	if pong, err := rdb.Ping(context.Background()).Result(); err != nil {
		fmt.Printf(" Redis connection failed: %v\n", err)
	} else {
		fmt.Printf("Redis connected: %s\n", pong)
	}

	return rdb
}
