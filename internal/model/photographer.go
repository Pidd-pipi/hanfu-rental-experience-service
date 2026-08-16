package model

import "time"

// Photographer is a photography service provider available for rental orders.
type Photographer struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Style       string    `gorm:"size:64" json:"style"`
	PricePerDay float64   `gorm:"not null" json:"price_per_day"`
	CreatedAt   time.Time `json:"created_at"`
}
