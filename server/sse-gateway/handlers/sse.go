package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sse/utils"
)

type ClipResponse struct {
	ID        string    `json:"id"`
	InBlob    bool      `json:"in_blob"`
	Content   string    `json:"content,omitempty"`
	BlobUrl   string    `json:"blob_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) SSEHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := utils.ExtractUserIDFromJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check consistent hash ownership
		owner := utils.HashRing.LocateKey([]byte(userID)).String()
		if owner != utils.GatewayAddr {
			http.Error(
				w,
				fmt.Sprintf("User %s belongs to %s, not %s", userID, owner, utils.GatewayAddr),
				http.StatusForbidden,
			)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// SSE setup
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		channelName := fmt.Sprintf("clips:user:%s", userID)
		pubsub := h.Redis.Subscribe(context.Background(), channelName)
		defer pubsub.Close()

		fmt.Printf("[SSE] User %s connected on %s\n", userID, utils.GatewayAddr)

		sendSSE(w, flusher, "connected", map[string]string{
			"message": fmt.Sprintf("Connected to %s", channelName),
		})

		ch := pubsub.Channel()
		for {
			select {
			case msg := <-ch:
				if msg == nil {
					return
				}

				var clip ClipResponse
				if err := json.Unmarshal([]byte(msg.Payload), &clip); err != nil {
					fmt.Println("[SSE] Invalid clip payload:", err)
					continue
				}

				// always send event type as new_clip
				sendSSE(w, flusher, "new_clip", clip)

			case <-r.Context().Done():
				fmt.Printf("[SSE] Client disconnected: user %s\n", userID)
				return
			}
		}
	}
}

func sendSSE(w http.ResponseWriter, flusher http.Flusher, eventType string, data any) {
	id := time.Now().UnixMilli()
	jsonData, _ := json.Marshal(data)
	_, _ = fmt.Fprintf(w, "id: %d\n", id)
	_, _ = fmt.Fprintf(w, "event: %s\n", eventType)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()
}
