package initializers

import (
	"os"

	"api-server/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	connUrl := os.Getenv("DATABASE_URL")
	Database, _ := gorm.Open(postgres.Open(connUrl), &gorm.Config{})

	return Database
}

func InitializeDatabaseAndMigrate() *gorm.DB {
	pg := ConnectDatabase()
	_ = pg.AutoMigrate(
		&models.User{},
		&models.Clip{},
		&models.ClipBlobMetadata{},
		&models.Outbox{},
	)
	return pg
}
