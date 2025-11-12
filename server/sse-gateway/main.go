package main

import (
	"fmt"
	"net/http"
	"os"

	"sse/handlers"
	"sse/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	redisClient := utils.ConnectRedisClient()
	utils.InitConsistentHashingRing()
	handler := handlers.CreateHandler(redisClient)

	port := os.Getenv("PORT")
	if port == "" {
		panic("Port not found for SSE gateway in env file")
	}

	http.HandleFunc("/events", handler.SSEHandler())

	fmt.Printf("[SSE] Gateway %s running on port %s\n", utils.GatewayAddr, port)
	fmt.Printf("[SSE] Connected to Redis and ready for clients...\n")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
