package clips

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"api-server/internal/models"
	apiUtils "api-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type ClipHandler struct {
	Service *ClipService
}

func NewHandler(s *ClipService) *ClipHandler {
	return &ClipHandler{
		Service: s,
	}
}

// Init Blob Upload

func (h *ClipHandler) InitBlobUpload(c *gin.Context) {
	var req models.BlobInitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	fileID, objectKey, uploadURL, err := h.Service.InitBlobUpload(req.MimeType)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to initialize upload"})
		return
	}

	c.JSON(200, gin.H{
		"id":         fileID,
		"object_key": objectKey,
		"upload_url": uploadURL,
	})
}

// Create Text Clip

func (h *ClipHandler) CreateTextClip(c *gin.Context) {
	var req models.CreateTextClipRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	userID := apiUtils.GetUserIdFromContext(c)

	clip, err := h.Service.CreateTextClip(userID, req.Content)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create clip"})
		return
	}

	c.JSON(201, gin.H{
		"id":        clip.ID,
		"inBlob":    clip.InBlob,
		"content":   clip.Content,
		"createdAt": clip.CreatedAt,
	})
}

// Create Blob Clip

func (h *ClipHandler) CreateBlobClip(c *gin.Context) {
	var req models.CreateBlobClipRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	userID := apiUtils.GetUserIdFromContext(c)

	clip, err := h.Service.CreateBlobClip(req.FileID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "blob not found") {
			c.JSON(400, gin.H{"error": "file not uploaded"})
			return
		}
		c.JSON(500, gin.H{"error": "failed to create blob clip"})
		return
	}

	c.JSON(201, gin.H{
		"id":        clip.ID,
		"inBlob":    true,
		"createdAt": clip.CreatedAt,
	})
}

// Get Clips

func (h *ClipHandler) GetClips(c *gin.Context) {
	before := c.Query("before")
	userID := apiUtils.GetUserIdFromContext(c)

	result, err := h.Service.GetClips(before, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to load clips"})
		return
	}

	c.JSON(200, result)
}

func (h *ClipHandler) GetLatestClip(c *gin.Context) {
	userID := apiUtils.GetUserIdFromContext(c)
	key := fmt.Sprintf("clips:latest:user:%d", userID)
	ctx := context.Background()

	// Check Redis
	val, err := h.Service.Cache.Get(ctx, key).Result()
	if err == nil {
		var clip models.ClipResponse
		_ = json.Unmarshal([]byte(val), &clip)

		c.JSON(200, clip)
		return
	}

	// Cache miss, fetch from DB and store in redis
	clip, err := h.Service.GetLatestClipFromDB(userID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to fetch latest clip"})
	}

	clipByte, _ := json.Marshal(&clip)
	_ = h.Service.Cache.Set(ctx, key, clipByte, 5*time.Minute)

	c.JSON(200, clip)
}
