package auth

import (
	"api-server/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *AuthService
}

func NewHandler(s *AuthService) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.Service.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := h.Service.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
