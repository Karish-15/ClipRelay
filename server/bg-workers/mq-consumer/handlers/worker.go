package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"consumer/initializers"
	"consumer/models"

	"github.com/minio/minio-go/v7"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type ConsumerHandler struct {
	Consumer *initializers.RabbitConsumer
	Redis    *redis.Client
	Blob     *minio.Client
}

func CreateHandler(Consumer *initializers.RabbitConsumer, Redis *redis.Client, Blob *minio.Client) *ConsumerHandler {
	return &ConsumerHandler{
		Consumer: Consumer,
		Redis:    Redis,
		Blob:     Blob,
	}
}

func (h *ConsumerHandler) ProcessWithRetries(d amqp091.Delivery, workerID int) {
	attempt := 0
	for {
		attempt++
		err := h.HandleEvent(d)
		if err == nil {
			if errAck := d.Ack(true); errAck != nil {
				log.Printf("[w%d] ack error: %v", workerID, errAck)
			}
			return
		}

		log.Printf("[w%d] processing error (attempt %d): %v", workerID, attempt, err)

		if attempt < 3 {
			backoff := time.Duration(100*(1<<attempt)) * time.Millisecond // exponential-ish
			time.Sleep(backoff)
			continue
		}

		// Reorder to queue, works because we maintain last-write-wins policy for redis
		if errNack := d.Nack(false, true); errNack != nil {
			log.Printf("[w%d] nack error: %v", workerID, errNack)
		}
		return
	}
}

func (h *ConsumerHandler) HandleEvent(d amqp091.Delivery) error {
	ctx := context.Background()

	var event models.OutboxEvent
	if err := json.Unmarshal(d.Body, &event); err != nil {
		log.Printf("Invalid Clip Event: %v\nevent:%v", err, d.Body)
		return nil
	}

	var payload models.ClipEventPayload
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		log.Printf("Invalid Clip Event Payload: %v\npayload:%v", err, event.Payload)
		return nil
	}

	var blobURL string
	if payload.InBlob && payload.Bucket != "" && payload.ObjectKey != "" {
		url, err := h.FetchBlobGetPresignedUrl(payload.ObjectKey, payload.Bucket)
		if err != nil {
			log.Printf("Failed to fetch blob URL for object ID: %v", payload.ObjectKey)
		}
		blobURL = url
	}

	clip := models.ClipResponse{
		ID:        payload.ID,
		InBlob:    payload.InBlob,
		Content:   payload.Content,
		BlobUrl:   blobURL,
		CreatedAt: payload.CreatedAt,
	}

	if err := h.RedisUpdateLatestAndPublish(&clip, payload.UserID, ctx); err != nil {
		log.Printf("Error while updating and publishing clip to redis: %v", err)
	}
	return nil
}
