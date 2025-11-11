package clips

import (
	"api-server/internal/middlewares"
	"api-server/internal/modules/outbox"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, blob *minio.Client, bus *outbox.Bus, r *gin.Engine) {
	service := NewService(db, blob, bus)
	handler := NewHandler(service)

	clips := r.Group("/clips")
	{
		clips.Use(middlewares.JWTMiddleware())
		clips.GET("", handler.GetClips)
		clips.POST("/text", handler.CreateTextClip)

		blobGroup := clips.Group("/blob")
		{
			blobGroup.POST("/init", handler.InitBlobUpload)
			blobGroup.POST("", handler.CreateBlobClip)
		}
	}
}
