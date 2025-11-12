package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"consumer/handlers"
	"consumer/initializers"

	"github.com/rabbitmq/amqp091-go"
)

const (
	workerFactor  = 2    // workers = CPU * workerFactor
	jobBufferSize = 1000 // buffered channel size for incoming deliveries
)

func main() {
	redisClient := initializers.ConnectRedisClient()
	rabbitConsumer, errRabbit := initializers.NewRabbitConsumer()
	if errRabbit != nil {
		log.Printf("error creating rabbit consumer: %v", errRabbit)
		panic("Rabbit init failed")
	}
	minIOClient := initializers.CreateAndInitMinIO()
	handler := handlers.CreateHandler(rabbitConsumer, redisClient, minIOClient)

	workers := runtime.NumCPU() * workerFactor
	if workers < 2 {
		workers = 2
	}

	if err := handler.Consumer.Channel.Qos(workers, 0, false); err != nil {
		log.Fatalf("set qos: %v", err)
	}

	msgs, _ := handler.Consumer.Channel.Consume(
		handler.Consumer.Queue.Name,
		"",
		false, // autoAck=false (manual ack)
		false,
		false,
		false,
		nil,
	)

	// Start workers
	jobs := make(chan amqp091.Delivery, jobBufferSize)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for d := range jobs {
				handler.ProcessWithRetries(d, workerID)
			}
		}(i)
	}

	// Start consuming from MQ, add deliveries to jobs
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					// channel closed
					close(jobs)
					return
				}
				jobs <- d
			case <-ctx.Done():
				// stop pumping
				close(jobs)
				return
			}
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	log.Println("shutdown signal received, stopping consumer...")
	cancel()
	// Wait for workers to finish current jobs
	wg.Wait()
	log.Println("workers finished, exiting")
}
