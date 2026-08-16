// Package model defines the GORM entity structures for hanfu-rental.
package model

import "time"

// User represents a customer or admin.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Phone        string    `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	PasswordHash string    `gorm:"size:100;not null" json:"-"`
	Nickname     string    `gorm:"size:32;not null" json:"nickname"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Role         string    `gorm:"size:16;not null;default:customer" json:"role"`
	Deposit      float64   `gorm:"not null;default:0" json:"deposit"`
	CreatedAt    time.Time `json:"created_at"`
}
