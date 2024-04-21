package responses

import "time"

type UserResponse struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt   *time.Time          `json:"deleted_at"`
	Email       string              `json:"email"`
	UserProfile UserProfileResponse `json:"user_profile"`
	LikedTags   []TagsResponse      `json:"liked_tags"`
}
