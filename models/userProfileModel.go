package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model

	DisplayName *string
	FirstName   *string
	MiddleName  *string
	LastName    *string
	Age         *uint8
	Birthday    *time.Time
	UserId      uint64
	User        User
}
