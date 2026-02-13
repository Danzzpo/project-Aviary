package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string    `gorm:"unique;not null" json:"username"`
	Email             string    `gorm:"unique;not null" json:"email"`
	PasswordHash      string    `json:"-"` 
	Role              string    `json:"role"`
	ProfilePic        string    `json:"profile_pic"` 
	LastProfileUpdate time.Time `json:"last_profile_update"` 
}