package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"consumer/models"

	"github.com/redis/go-redis/v9"
)

func (h *ConsumerHandler) FetchBlobGetPresignedUrl(objectID string, bucket string) (string, error) {
	url, err := h.Blob.PresignedGetObject(
		context.Background(),
		bucket,
		objectID,
		30*time.Minute,
		nil,
	)

	return url.String(), err
}

func (h *ConsumerHandler) RedisFetchClip(key string, ctx context.Context) (*models.ClipResponse, error) {
	var clip models.ClipResponse
	resp, err := h.Redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis key doesnt exist, key: %v", key)
			return nil, nil
		}
		log.Printf("Redis GET Failed: %v", key)
		return nil, err
	}

	_ = json.Unmarshal([]byte(resp), &clip)
	return &clip, nil
}

func (h *ConsumerHandler) RedisUpdateLatestAndPublish(clip *models.ClipResponse, userID int, ctx context.Context) error {
	key := fmt.Sprintf("clip:latest:user:%d", userID)

	existing, err := h.RedisFetchClip(key, ctx)
	if err != nil {
		return err
	}

	if existing != nil {
		if clip.CreatedAt.Before(existing.CreatedAt) {
			log.Printf("Ignoring Older clip for user ID: %v", userID)
			return nil
		}
	}

	// Update latest key
	data, _ := json.Marshal(clip)
	if err = h.Redis.Set(ctx, key, data, 5*time.Minute).Err(); err != nil {
		log.Printf("Error setting clip for user ID: %v. error: %v", userID, err)
	}
	log.Printf("Redis clip updated for user %d -> key=%s", userID, key)

	// Publish key to Redis pubsub
	channel := fmt.Sprintf("clip:user:%d", userID)
	if err = h.Redis.Publish(ctx, channel, data).Err(); err != nil {
		log.Printf("Error publishing clip for user ID: %v. error: %v", userID, err)
	}
	log.Printf("Redis clip published for user ID: %d", userID)

	return nil
}
