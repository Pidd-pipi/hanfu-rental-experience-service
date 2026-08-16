package dto

import "time"

// CreateActivityRequest is the payload for publishing an activity.
type CreateActivityRequest struct {
	Title           string    `json:"title" binding:"required,min=2,max=64"`
	Category        string    `json:"category" binding:"required"`
	StartTime       time.Time `json:"start_time" binding:"required"`
	Location        string    `json:"location"`
	Flow            string    `json:"flow"`
	DressCode       string    `json:"dress_code"`
	Fee             float64   `json:"fee" binding:"gte=0"`
	MaxParticipants int       `json:"max_participants" binding:"required,gte=2"`
}
