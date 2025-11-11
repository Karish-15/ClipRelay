package main

import (
	"context"

	"api-server/internal/initializers"
	"api-server/internal/modules"
	"api-server/internal/modules/outbox"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := initializers.InitializeDatabaseAndMigrate()
	blob := initializers.CreateAndInitMinIO()
	bus := outbox.NewBus(db)
	cache := initializers.ConnectRedisClient()
	consistentRing := initializers.InitConsistentHashingRing()

	bus.Start(ctx)
	r := gin.Default()

	modules.RegisterAll(r, db, blob, cache, bus, consistentRing)

	if err := r.Run(":8080"); err != nil {
		panic("Error! Failed to start application.")
	}
}
