package modules

import (
	"api-server/internal/modules/auth"
	"api-server/internal/modules/clips"
	"api-server/internal/modules/users"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func RegisterAll(r *gin.Engine, db *gorm.DB, blob *minio.Client) {
	users.Register(db, r)
	clips.Register(db, blob, r)
	auth.Register(db, r)
}
