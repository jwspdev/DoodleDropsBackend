package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey"`
	AuthorID  uint64
	Author    *User
	Content   *string
	LikedBy   *[]User
	Comments  *[]Comment
	Tags      *[]Tag
	CreatedAt time.Time
}
