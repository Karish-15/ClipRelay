package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	pollMS := getInt("OUTBOX_WORKER_POLL_MS", 200)
	batchSize := getInt("OUTBOX_WORKER_BATCH_SIZE", 100)
	leaseTimeout := time.Duration(getInt("OUTBOX_WORKER_LEASE_TIMEOUT_SECONDS", 60)) * time.Second

	repo := NewOutboxRepo(db, leaseTimeout)
	publisher, err := NewRabbitPublisher()
	if err != nil {
		log.Fatalf("rabbitmq failed: %v", err)
	}
	defer publisher.Close()

	workerID := "worker-1" // change in docker deployments if needed

	log.Println("Outbox Worker started...")

	for {
		events, err := repo.FetchAndClaim(ctx, workerID, batchSize)
		if err != nil {
			log.Printf("fetch error: %v", err)
			sleep(pollMS)
			continue
		}

		if len(events) == 0 {
			sleep(pollMS)
			continue
		}

		ids := extractIDs(events)
		publishOK := true

		for _, evt := range events {
			if err := publisher.Publish(evt); err != nil {
				publishOK = false
				log.Printf("publish failed for event=%d: %v", evt.ID, err)
				break
			}
		}

		if publishOK {
			if err := repo.MarkProcessed(ctx, ids); err != nil {
				log.Printf("mark processed error: %v", err)
			}
		} else {
			if err := repo.Release(ctx, ids); err != nil {
				log.Printf("release events error: %v", err)
			}
		}

		log.Printf("published %d events to queue", len(ids))
		sleep(pollMS)
	}
}

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func getInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

func extractIDs(events []OutboxEvent) []int64 {
	out := make([]int64, len(events))
	for i, evt := range events {
		out[i] = evt.ID
	}
	return out
}
