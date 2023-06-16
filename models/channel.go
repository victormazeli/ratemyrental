package models

type Channel struct {
	ID          uint `gorm:"primaryKey"`
	SenderID    uint `json:"sender_id"`
	RecipientID uint `json:"recipient_id"`
}
