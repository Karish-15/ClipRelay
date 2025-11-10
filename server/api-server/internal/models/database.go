package models

import "time"

type User struct {
	ID           int       `gorm:"primaryKey;index"`
	Username     string    `gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string    `gorm:"size:255"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type Clip struct {
	ID     string `gorm:"primaryKey;type:uuid"`
	UserID int    `gorm:"index;not null"`

	// Present in Blob?
	InBlob bool `gorm:"default:false"`

	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationship to User
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type ClipBlobMetadata struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	ClipID    string    `gorm:"uniqueIndex;not null"`
	Bucket    string    `gorm:"size:100;not null"`
	ObjectKey string    `gorm:"size:255;not null"`
	Status    string    `gorm:"size:20;default:'PENDING'"` // NEW
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Outbox struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	EventType string    `gorm:"size:100;not null"`   // e.g., "clip.created"
	Payload   string    `gorm:"type:jsonb;not null"` // serialized event payload
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Processed bool      `gorm:"default:false"` // worker flips this after publishing
}
