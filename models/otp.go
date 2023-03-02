package models

import (
	"gorm.io/gorm"
	"time"
)

type Otp struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Email      string         `json:"email" gorm:"type:varchar(255);not null""`
	Purpose    string         `json:"purpose" gorm:"type:varchar(255);not null""`
	Status     uint8          `json:"status"`
	Otp        string         `gorm:"type:varchar(255);not null""`
	ExpiryDate time.Time      `json:"expiry_date"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
