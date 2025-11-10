package users

import (
	"api-server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	Service *UserService
}

func NewHandler(s *UserService) *Handler {
	return &Handler{Service: s}
}

// Create User
func (h *Handler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.Service.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(201, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Get User by ID
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(500, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Update User
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.Service.UpdateUser(id, req.Username, req.Password)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(500, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Delete User by Username
func (h *Handler) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	err := h.Service.DeleteUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(500, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}
