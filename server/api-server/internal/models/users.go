package models

// Request Models
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Response
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
