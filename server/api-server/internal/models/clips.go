package models

import "time"

type CreateTextClipRequest struct {
	Content string `json:"content" binding:"required"`
}

type BlobInitRequest struct {
	MimeType string `json:"mime_type" binding:"required"`
}

type CreateBlobClipRequest struct {
	FileID string `json:"id" binding:"required"`
}

// RESPONSE

type ClipResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content,omitempty"`
	InBlob    bool      `json:"in_blob"`
	CreatedAt time.Time `json:"created_at"`
}
