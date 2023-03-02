package models

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

type User struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	FullName   null.String `json:"full_name" gorm:"type:varchar(255)"`
	Email      string      `json:"email" gorm:"type:varchar(255);unique;not null""`
	Password   string      `json:"password" gorm:"type:varchar(255);not null""`
	Avatar     null.String `json:"avatar" gorm:"type:varchar(255)"`
	Address    null.String `json:"address" gorm:"type:varchar(255)"`
	PostalCode null.String `json:"postal_code" gorm:"type:varchar(255)"`
	Location   null.String `json:"location" gorm:"type:text"`
	Status     uint8       `json:"status"` // 0 = not verified; 1 = verified; 2 = suspended; 3 = deleted; 4 = locked;
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  null.Time   `json:"deleted_at"`
}
