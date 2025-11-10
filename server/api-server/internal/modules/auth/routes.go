package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, r *gin.Engine) {
	service := NewService(db)
	handler := NewHandler(service)

	r.POST("/login", handler.Login)
}
