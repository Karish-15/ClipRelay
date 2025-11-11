package models

import "time"

type ClipCreatedPayload struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	InBlob    bool      `json:"in_blob"`
	Content   string    `json:"content,omitempty"`
	Bucket    string    `json:"bucket,omitempty"`
	ObjectKey string    `json:"object_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
