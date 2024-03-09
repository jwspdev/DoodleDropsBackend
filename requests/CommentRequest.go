package requests

type CreateCommentRequest struct {
	PostId  uint64 `json:"post_id"`
	Content string `json:"content"`
}
