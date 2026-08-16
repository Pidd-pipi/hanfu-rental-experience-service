package model

import "time"

// Hanfu is a hanfu garment in the rental catalog.
type Hanfu struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Dynasty     string    `gorm:"size:16;index;not null" json:"dynasty"`
	Form        string    `gorm:"size:32" json:"form"`
	Color       string    `gorm:"size:32" json:"color"`
	Size        string    `gorm:"size:16" json:"size"`
	PricePerDay float64   `gorm:"not null" json:"price_per_day"`
	Images      string    `gorm:"type:text" json:"images"`
	Stock       int       `gorm:"not null;default:1" json:"stock"`
	Status      string    `gorm:"size:16;index;not null;default:available" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
