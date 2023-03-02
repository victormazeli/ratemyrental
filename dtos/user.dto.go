package dtos

import (
	"time"
)

type UserDTO struct {
	Id         uint      `json:"id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	Location   string    `json:"location"`
	Status     uint8     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
