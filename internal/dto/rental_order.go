package dto

import "time"

// CreateRentalOrderRequest is the combined rental order payload.
type CreateRentalOrderRequest struct {
	HanfuID         uint      `json:"hanfu_id" binding:"required"`
	MakeupPackageID *uint     `json:"makeup_package_id"`
	PhotographerID  *uint     `json:"photographer_id"`
	StartDate       time.Time `json:"start_date" binding:"required"`
	EndDate         time.Time `json:"end_date" binding:"required"`
}
