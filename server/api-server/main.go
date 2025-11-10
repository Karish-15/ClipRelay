package main

import (
	"api-server/internal/initializers"
	"api-server/internal/modules"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db := initializers.InitializeDatabaseAndMigrate()
	blob := initializers.CreateAndInitMinIO()

	r := gin.Default()

	modules.RegisterAll(r, db, blob)

	if err := r.Run(":8080"); err != nil {
		panic("Error! Failed to start application.")
	}
}
