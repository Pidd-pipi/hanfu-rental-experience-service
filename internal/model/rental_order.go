package model

import "time"

// RentalOrder combines hanfu rental with makeup and photography options.
type RentalOrder struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"index;not null" json:"user_id"`
	HanfuID          uint       `gorm:"index;not null" json:"hanfu_id"`
	MakeupPackageID  *uint      `json:"makeup_package_id"`
	PhotographerID   *uint      `json:"photographer_id"`
	StartDate        time.Time  `gorm:"not null" json:"start_date"`
	EndDate          time.Time  `gorm:"not null" json:"end_date"`
	DepositAmount    float64    `gorm:"not null" json:"deposit_amount"`
	RentalFee        float64    `gorm:"not null" json:"rental_fee"`
	Status           string     `gorm:"size:16;index;not null;default:pending" json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
}
