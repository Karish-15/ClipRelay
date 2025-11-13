package initializers

import (
	"os"
	"time"

	"api-server/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	connUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(connUrl), &gorm.Config{
		// DISABLE prepared statements (they slow down high-throughput inserts)
		PrepareStmt: false,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return db
}

func InitializeDatabaseAndMigrate() *gorm.DB {
	db := ConnectDatabase()

	// Run AutoMigrate ONCE at startup, not during load.
	_ = db.AutoMigrate(
		&models.User{},
		&models.Clip{},
		&models.ClipBlobMetadata{},
		&models.Outbox{},
	)

	return db
}
