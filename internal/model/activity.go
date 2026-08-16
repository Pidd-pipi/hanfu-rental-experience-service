package model

import "time"

// ActivityRegistration is a customer's signup for a cultural activity.
type ActivityRegistration struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ActivityID uint      `gorm:"uniqueIndex:uniq_activity_registration,priority:1;not null" json:"activity_id"`
	UserID     uint      `gorm:"uniqueIndex:uniq_activity_registration,priority:2;not null" json:"user_id"`
	Fee        float64   `gorm:"not null" json:"fee"`
	CreatedAt  time.Time `json:"created_at"`
}

// Activity is a cultural experience event published by the store.
type Activity struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Title           string    `gorm:"size:64;not null" json:"title"`
	Category        string    `gorm:"size:32;index;not null" json:"category"`
	StartTime       time.Time `gorm:"not null" json:"start_time"`
	Location        string    `gorm:"size:128" json:"location"`
	Flow            string    `gorm:"type:text" json:"flow"`
	DressCode       string    `gorm:"size:128" json:"dress_code"`
	Fee             float64   `gorm:"not null" json:"fee"`
	MaxParticipants int       `gorm:"not null" json:"max_participants"`
	Status          string    `gorm:"size:16;not null;default:open" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}
