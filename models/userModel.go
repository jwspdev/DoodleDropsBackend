package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID    uint64 `gorm:"primaryKey"`
	Email string `gorm:"unique"`
	// Test     string
	Password    string
	UserProfile *UserProfile
}
