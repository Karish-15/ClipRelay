package modules

import (
	"api-server/internal/modules/auth"
	"api-server/internal/modules/clips"
	"api-server/internal/modules/gatewayhash"
	"api-server/internal/modules/outbox"
	"api-server/internal/modules/users"

	"github.com/buraksezer/consistent"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterAll(r *gin.Engine, db *gorm.DB, blob *minio.Client, cache *redis.Client, bus *outbox.Bus, ring *consistent.Consistent) {
	users.Register(db, r)
	clips.Register(db, blob, cache, bus, r)
	auth.Register(db, r)
	gatewayhash.Register(ring, r)
}
