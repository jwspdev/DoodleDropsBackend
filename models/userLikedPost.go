package models

import "gorm.io/gorm"

type UserLikedPost struct {
	gorm.Model
	UserId uint64
	PostId uint64
}
