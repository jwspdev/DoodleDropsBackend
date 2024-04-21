package requests

type CreatePostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
