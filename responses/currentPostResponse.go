package responses

type CurrentPostResponse struct {
	//TODO add image
	IsEditable  bool              `json:"is_editable"`
	Description string            `json:"description"`
	ID          uint64            `json:"id"`
	CreatedBy   CreatedByResponse `json:"user_details"`
}

type PostsTags struct {
}

type CreatedByResponse struct {
	//TODO add image of user
	UserID      uint64 `json:"user_id"`
	DisplayName string `json:"display_name"`
}

type PostCommentsResponse struct {
}
