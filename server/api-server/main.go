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

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	modules.RegisterAll(r, db, blob, cache, bus, consistentRing)

	if err := r.Run(":8080"); err != nil {
		panic("Error! Failed to start application.")
	}
}
