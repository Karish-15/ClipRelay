package models

import "time"

type OutboxEvent struct {
	ID        int64
	EventType string
	Payload   string
}

type ClipEventPayload struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	InBlob    bool      `json:"in_blob"`
	Content   string    `json:"content,omitempty"`
	Bucket    string    `json:"bucket,omitempty"`
	ObjectKey string    `json:"object_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ClipResponse struct {
	ID        string    `json:"id"`
	InBlob    bool      `json:"in_blob"`
	Content   string    `json:"content,omitempty"`
	BlobUrl   string    `json:"blob_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
