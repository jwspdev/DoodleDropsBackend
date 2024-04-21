package responses

import "time"

type CommentResponse struct {
	ID        uint64            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	CreatedBy CreatedByResponse `json:"created_by"`
	Content   string            `json:"content"`
}
