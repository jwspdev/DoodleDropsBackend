package responses

type DefaultPostResponse struct {
	//TODO ADD IMAGE
	// CreatedAt         time.Time      `json:"created_at"`
	// UpdatedAt         time.Time      `json:"updated_at"`
	ID uint64 `json:"id"`
	// Description       string         `json:"description"`
	Title             string         `json:"title"`
	Tags              []TagsResponse `json:"tags"`
	PostCounters      PostCounters   `json:"post_counters"`
	CreatedByResponse `json:"author"`
}

type CreatedByResponse struct {
	//TODO add image of user
	UserID      uint64 `json:"user_id"`
	DisplayName string `json:"display_name"`
}

type PostCounters struct {
	CommentCount   int `json:"comment_count"`
	UserLikedCount int `json:"user_liked_count"`
}
