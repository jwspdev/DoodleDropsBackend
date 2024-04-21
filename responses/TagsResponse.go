package responses

type TagsResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TagType     string `json:"tag_type"`
}
