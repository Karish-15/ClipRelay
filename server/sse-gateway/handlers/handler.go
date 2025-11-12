package handlers

import (
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Redis *redis.Client
}

func CreateHandler(Redis *redis.Client) *Handler {
	return &Handler{
		Redis: Redis,
	}
}
