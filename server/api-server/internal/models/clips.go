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
	InBlob    bool      `json:"in_blob"`
	Content   string    `json:"content,omitempty"`
	BlobUrl   string    `json:"blob_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
