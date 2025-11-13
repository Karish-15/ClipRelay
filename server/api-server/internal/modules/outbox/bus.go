package outbox

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"api-server/internal/initializers"
	"api-server/internal/models"

	"gorm.io/gorm"
)

type Event struct {
	EventType string
	Payload   string
}

type Bus struct {
	channel   chan Event
	database  *gorm.DB
	batchSize int
}

func NewBus(db *gorm.DB) *Bus {
	qSize := initializers.GetEnvInt("OUTBOX_WORKER_QUEUE_SIZE", 10000)
	batchSize := initializers.GetEnvInt("OUTBOX_WORKER_BATCH_SIZE", 100)

	return &Bus{
		channel:   make(chan Event, qSize),
		database:  db,
		batchSize: batchSize,
	}
}

// -------- Async Enqueue
func (b *Bus) Enqueue(eventType string, payload any) {
	data, _ := json.Marshal(payload)

	event := Event{
		EventType: eventType,
		Payload:   string(data),
	}

	select {
	case b.channel <- event:
		// queued
	default:
		// fallback -> sync update
		_ = b.insertOne(event)
	}
}

func (b *Bus) insertOne(e Event) error {
	return b.database.Create(&models.Outbox{
		EventType: e.EventType,
		Payload:   e.Payload,
	}).Error
}

// -------- Worker goroutine
func (b *Bus) Start(ctx context.Context) {
	go func() {
		batch := make([]models.Outbox, 0, b.batchSize)
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				b.flush(batch)
				return

			case e := <-b.channel:
				batch = append(batch, models.Outbox{
					EventType: e.EventType,
					Payload:   e.Payload,
				})

				// batch full → flush
				if len(batch) >= b.batchSize {
					b.flush(batch)
					batch = batch[:0]
				}

			case <-ticker.C:
				// timed flush
				b.flush(batch)
				batch = batch[:0]
			}
		}
	}()
}

// -------- Bulk Insert
func (b *Bus) flush(batch []models.Outbox) {
	if len(batch) == 0 {
		return
	}
	log.Printf("Flushing %d rows to Outbox table", len(batch))
	if err := b.database.CreateInBatches(batch, b.batchSize).Error; err != nil {
		log.Println("[outbox] flush error:", err)
	}
}
