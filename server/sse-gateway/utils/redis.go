package utils

import (
	"context"
	"os"

	redis "github.com/redis/go-redis/v9"
)

func ConnectRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:       0,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_ = client.Ping(context.Background()).Err()
	return client
}
