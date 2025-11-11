package gatewayhash

import (
	"api-server/internal/middlewares"

	"github.com/buraksezer/consistent"
	"github.com/gin-gonic/gin"
)

func Register(Ring *consistent.Consistent, r *gin.Engine) {
	handler := NewHandler(Ring)

	gatewayGroup := r.Group("/sse")
	{
		gatewayGroup.Use(middlewares.JWTMiddleware())
		gatewayGroup.GET("", handler.ResolveGateway)
	}
}
