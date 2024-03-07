package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey"`
	AuthorID uint64
	Author   User
	PostId   *uint64
	Post     *Post
	Content  string
	ParentID *uint64
	// Parent   *Comment `gorm:"foreignkey:ParentID"`
	// Replies  []*Comment
}
