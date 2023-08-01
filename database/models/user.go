package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type User struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	FullName      null.String `json:"full_name" gorm:"type:varchar(255)"`
	Email         string      `json:"email" gorm:"type:varchar(255);unique;not null""`
	Password      string      `json:"password" gorm:"type:varchar(255);not null""`
	Avatar        null.String `json:"avatar" gorm:"type:varchar(255)"`
	Latitude      float64     `json:"latitude"`
	Longitude     float64     `json:"longitude"`
	City          string      `json:"city" gorm:"type:varchar(255)"`
	Country       string      `json:"country" gorm:"type:varchar(255)"`
	PostalCode    null.String `json:"postal_code" gorm:"type:varchar(255)"`
	PhoneNumber   null.String `json:"phone_number" gorm:"type:varchar(255)"`
	Status        uint8       `json:"status"`                          // 0 = not verified; 1 = verified; 2 = suspended; 3 = deleted; 4 = locked;
	ChatStatus    uint8       `json:"chat_status"`                     // 0 = offline; 1 = online;
	UserType      uint8       `json:"user_type"`                       // 1 = regular-user; 2 = agent;
	ProfileStatus uint8       `json:"profile_status" gorm:"default:0"` // 1 or 0
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     null.Time   `json:"deleted_at"`
}
