package clips

import (
	"context"
	"fmt"
	"os"
	"time"

	"api-server/internal/constants"
	"api-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ClipService struct {
	DB   *gorm.DB
	Blob *minio.Client
}

func NewService(db *gorm.DB, blob *minio.Client) *ClipService {
	return &ClipService{
		DB:   db,
		Blob: blob,
	}
}

// Init Blob Upload

func (s *ClipService) InitBlobUpload(mime string) (string, string, string, error) {
	fileID := uuid.NewString()
	objectKey := "clips/" + fileID

	url, err := s.Blob.PresignedPutObject(
		context.Background(),
		os.Getenv("MINIO_BUCKET"),
		objectKey,
		10*time.Minute,
	)
	if err != nil {
		return "", "", "", err
	}

	meta := models.ClipBlobMetadata{
		ClipID:    fileID,
		Bucket:    os.Getenv("MINIO_BUCKET"),
		ObjectKey: objectKey,
		Status:    string(constants.PENDING),
	}

	if err := s.DB.Create(&meta).Error; err != nil {
		return "", "", "", err
	}

	return fileID, objectKey, url.String(), nil
}

// Create Text Clip

func (s *ClipService) CreateTextClip(userID int, content string) (*models.Clip, error) {
	clip := &models.Clip{
		ID:      uuid.NewString(),
		UserID:  userID,
		InBlob:  false,
		Content: content,
	}

	if err := s.DB.Create(clip).Error; err != nil {
		return nil, err
	}

	return clip, nil
}

// Update DB Entry once file uploaded

func (s *ClipService) CreateBlobClip(fileID string, userID int) (*models.Clip, error) {
	var meta models.ClipBlobMetadata
	if err := s.DB.First(&meta, "clip_id = ?", fileID).Error; err != nil {
		return nil, err
	}

	// Verify object exists in MinIO
	_, err := s.Blob.StatObject(
		context.Background(),
		meta.Bucket,
		meta.ObjectKey,
		minio.StatObjectOptions{},
	)
	if err != nil {
		// If object is missing or 0 byte → fail
		return nil, fmt.Errorf("blob not found or not uploaded: %w", err)
	}

	// Create clip record
	clip := &models.Clip{
		ID:      fileID,
		UserID:  userID,
		InBlob:  true,
		Content: "",
	}

	if err := s.DB.Create(clip).Error; err != nil {
		return nil, err
	}

	// update metadata status
	meta.Status = string(constants.UPLOADED)
	if err := s.DB.Save(&meta).Error; err != nil {
		return nil, err
	}

	return clip, nil
}

// Clips with pagination
func (s *ClipService) GetClips(before string, userID int) ([]gin.H, error) {
	query := s.DB.
		Model(&models.Clip{}).
		Order("created_at desc").
		Limit(5).
		Where("user_id = ?", userID)

	if before != "" {
		var ref models.Clip
		if err := s.DB.First(&ref, "id = ?", before).Error; err == nil {
			query = query.Where("created_at < ?", ref.CreatedAt)
		}
	}

	var clips []models.Clip
	if err := query.Find(&clips).Error; err != nil {
		return nil, err
	}

	result := make([]gin.H, 0, len(clips))

	for _, clip := range clips {
		item := gin.H{
			"id":         clip.ID,
			"in_blob":    clip.InBlob,
			"created_at": clip.CreatedAt,
		}

		if !clip.InBlob {
			item["content"] = clip.Content
			result = append(result, item)
			continue
		}

		// MANUAL METADATA LOOKUP
		var meta models.ClipBlobMetadata
		if err := s.DB.First(&meta, "clip_id = ?", clip.ID).Error; err != nil {
			item["error"] = "blob metadata missing"
			result = append(result, item)
			continue
		}

		url, err := s.Blob.PresignedGetObject(
			context.Background(),
			meta.Bucket,
			meta.ObjectKey,
			30*time.Minute,
			nil,
		)
		if err != nil {
			item["error"] = "failed to generate blob URL"
		} else {
			item["blob_url"] = url.String()
		}

		result = append(result, item)
	}

	return result, nil
}
