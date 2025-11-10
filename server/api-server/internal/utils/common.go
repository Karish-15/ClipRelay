package utils

import "github.com/gin-gonic/gin"

func GetUserIdFromContext(c *gin.Context) int {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(int)
	return userID
}
