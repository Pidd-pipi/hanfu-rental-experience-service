package util

import (
	"math"
	"time"
)

// RentalFeeResult is the computed rental fee payload.
type RentalFeeResult struct {
	Days        int     `json:"days"`
	Subtotal    float64 `json:"subtotal"`
	CardDiscount float64 `json:"card_discount"`
	MakeupFee   float64 `json:"makeup_fee"`
	PhotographerFee float64 `json:"photographer_fee"`
	Total       float64 `json:"total"`
}

// CalcRentalDays returns the inclusive day count between two dates.
func CalcRentalDays(start, end time.Time) int {
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	days := int(math.Round(e.Sub(s).Hours()/24))
	if days < 1 {
		days = 1
	}
	return days
}

// CalcRentalFee computes subtotal, member discount and total rental fee.
// Member card holders get a 15% discount; annual card holders get 25%.
func CalcRentalFee(pricePerDay float64, start, end time.Time, hasCard bool, cardType string, makeupFee, photographerFee float64) RentalFeeResult {
	days := CalcRentalDays(start, end)
	subtotal := pricePerDay * float64(days)
	discountRate := 0.0
	if hasCard {
		discountRate = 0.15
		if cardType == "month" {
			discountRate = 0.25
		}
	}
	discount := subtotal * discountRate
	total := subtotal - discount + makeupFee + photographerFee
	return RentalFeeResult{
		Days: days, Subtotal: round2(subtotal), CardDiscount: round2(discount),
		MakeupFee: round2(makeupFee), PhotographerFee: round2(photographerFee),
		Total: round2(total),
	}
}

// CalcDeposit returns the deposit amount (one day's rent or 200 min).
func CalcDeposit(pricePerDay float64) float64 {
	d := pricePerDay
	if d < 200 {
		d = 200
	}
	return round2(d)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
