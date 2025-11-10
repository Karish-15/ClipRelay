package users

import (
	"api-server/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, r *gin.Engine) {
	service := NewService(db)
	handler := NewHandler(service)

	users := r.Group("/users")
	{
		users.POST("", handler.CreateUser)

		protectedRoutes := users.Group("")
		{
			protectedRoutes.Use(middlewares.JWTMiddleware())
			protectedRoutes.GET("/:id", handler.GetUser)
			protectedRoutes.PUT("/:id", handler.UpdateUser)
			protectedRoutes.DELETE("/:username", handler.DeleteUser)
		}

	}
}
