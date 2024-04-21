package responses

import "time"

//TODO ADD MORE FIELDS
type CurrentPostResponse struct {
	//TODO add image
	Description string    `json:"description"`
	IsEditable  bool      `json:"is_editable"`
	CreatedAt   time.Time `json:"created_at"`
}

type PostsTags struct {
}

type PostCommentsResponse struct {
}
