package model

import "time"

// MemberCard is a rental membership card owned by a customer.
type MemberCard struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	CardType  string    `gorm:"size:16;not null" json:"card_type"`
	Status    string    `gorm:"size:16;not null;default:active" json:"status"`
	ExpireAt  time.Time `gorm:"not null" json:"expire_at"`
	CreatedAt time.Time `json:"created_at"`
}
