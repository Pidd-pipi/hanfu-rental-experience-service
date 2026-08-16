package model

import "time"

// MakeupPackage is a makeup styling bundle available with rental orders.
type MakeupPackage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}
