package util

import (
	"testing"
	"time"
)

func TestCalcRentalDays(t *testing.T) {
	start := time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		end  time.Time
		want int
	}{
		{name: "same day", end: time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC), want: 1},
		{name: "two days", end: time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC), want: 2},
		{name: "three days", end: time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC), want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcRentalDays(start, tt.end); got != tt.want {
				t.Fatalf("CalcRentalDays = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalcRentalFee(t *testing.T) {
	start := time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC) // 3 days
	tests := []struct {
		name    string
		hasCard bool
		cardType string
		wantTotal float64
	}{
		{name: "no card", hasCard: false, wantTotal: 300},
		{name: "month card 15%", hasCard: true, cardType: "month", wantTotal: 255},
		{name: "year card 25%", hasCard: true, cardType: "year", wantTotal: 225},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CalcRentalFee(100, start, end, tt.hasCard, tt.cardType, 0, 0)
			if res.Total != tt.wantTotal {
				t.Fatalf("total = %f, want %f", res.Total, tt.wantTotal)
			}
		})
	}
}

func TestCalcDeposit(t *testing.T) {
	if got := CalcDeposit(100); got != 200 {
		t.Fatalf("CalcDeposit(100) = %f, want 200", got)
	}
	if got := CalcDeposit(300); got != 600 {
		t.Fatalf("CalcDeposit(300) = %f, want 600", got)
	}
}
