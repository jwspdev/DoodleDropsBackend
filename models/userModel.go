package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID    uint64 `gorm:"primaryKey"`
	Email string `gorm:"unique"`
	// Test     string
	Password    string
	UserProfile *UserProfile
	Posts       []*Post    `gorm:"foreignKey:AuthorID"`
	Comments    []*Comment `gorm:"foreignKey:AuthorID"`
	LikedPosts  []*Post    `gorm:"many2many:user_liked_posts"`
	LikedTags   []*Tag     `gorm:"many2many:user_liked_tags"`
}
