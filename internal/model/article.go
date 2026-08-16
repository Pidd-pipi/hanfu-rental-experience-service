package model

import "time"

// Article is a hanfu culture science article published by the admin.
type Article struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:64;not null" json:"title"`
	Category    string    `gorm:"size:32" json:"category"`
	Content     string    `gorm:"type:longtext" json:"content"`
	CoverURL    string    `gorm:"size:255" json:"cover_url"`
	PublishedAt time.Time `json:"published_at"`
}
