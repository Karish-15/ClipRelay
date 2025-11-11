package outbox

import (
	"api-server/internal/constants"
	"api-server/internal/models"
)

func (bus *Bus) CreateClipEvent(clip models.Clip, meta *models.ClipBlobMetadata) {
	payload := models.ClipCreatedPayload{
		ID:        clip.ID,
		UserID:    clip.UserID,
		InBlob:    clip.InBlob,
		CreatedAt: clip.CreatedAt,
	}

	if clip.InBlob && meta != nil {
		payload.Bucket = meta.Bucket
		payload.ObjectKey = meta.ObjectKey
	} else {
		payload.Content = clip.Content
	}

	bus.Enqueue(constants.EventClipCreated, payload)
}
