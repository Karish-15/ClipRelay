package clips

import (
	"api-server/internal/middlewares"
	"api-server/internal/modules/outbox"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, blob *minio.Client, cache *redis.Client, bus *outbox.Bus, r *gin.Engine) {
	service := NewService(db, blob, cache, bus)
	handler := NewHandler(service)

	clips := r.Group("/clips")
	{
		clips.Use(middlewares.JWTMiddleware())
		clips.GET("", handler.GetClips)
		clips.POST("/text", handler.CreateTextClip)
		clips.GET("/latest", handler.GetLatestClip)

		blobGroup := clips.Group("/blob")
		{
			blobGroup.POST("/init", handler.InitBlobUpload)
			blobGroup.POST("", handler.CreateBlobClip)
		}
	}
}
