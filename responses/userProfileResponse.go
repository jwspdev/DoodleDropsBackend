package responses

import "time"

type UserProfileResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt   *time.Time `json:"deleted_at"`
	DisplayName *string    `json:"display_name"`
	FirstName   *string    `json:"first_name"`
	MiddleName  *string    `json:"middle_name"`
	LastName    *string    `json:"last_name"`
	Age         *uint8     `json:"age"`
	Birthday    *time.Time `json:"birthday"`
	UserID      *uint64    `json:"user_id"`
}
