package models

import "time"

type ChatMessage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Sender    uint      `json:"sender" gorm:"not null"`
	Recipient uint      `json:"recipient" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	ChannelID uint      `json:"channel_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
