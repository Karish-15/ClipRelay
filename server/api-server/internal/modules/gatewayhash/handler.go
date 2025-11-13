package gatewayhash

import (
	"fmt"

	"github.com/buraksezer/consistent"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Ring *consistent.Consistent
}

func NewHandler(Ring *consistent.Consistent) *Handler {
	return &Handler{
		Ring: Ring,
	}
}

func (h *Handler) ResolveGateway(c *gin.Context) {
	userIDStr := c.GetString("userID")

	fmt.Printf("HASH-KEY-DEBUG: '%s'\n", userIDStr)

	member := h.Ring.LocateKey([]byte(userIDStr))
	gateway := member.String()

	c.JSON(200, gin.H{
		"gateway": gateway,
	})
}
