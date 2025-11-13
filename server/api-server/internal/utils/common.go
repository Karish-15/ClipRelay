package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromContext(c *gin.Context) int {
	userIDVal, _ := c.Get("userID")
	userID, _ := strconv.Atoi(userIDVal.(string))
	return userID
}
