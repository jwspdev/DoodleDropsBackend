package models

import "gorm.io/gorm"

type UserLikedTags struct {
	gorm.Model
	UserId uint64
	PostId uint64
}
