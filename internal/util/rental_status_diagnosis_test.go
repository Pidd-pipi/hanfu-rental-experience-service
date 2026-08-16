package util

import (
	"testing"

	"github.com/lp/hanfu-rental/internal/constants"
)

func TestRentalStatusAndCardBoundary(t *testing.T) {
	if !constants.IsRentalOrderStatus(constants.RentalOrderStatusReturned) {
		t.Fatal("IsRentalOrderStatus(returned) should be true")
	}
	if got := RentalStatusText(constants.RentalOrderStatusPending); got != "待确认" {
		t.Fatalf("RentalStatusText(pending) = %s, want 待确认", got)
	}
	if got := constants.MemberCardTypeText(constants.MemberCardTypeYear); got != "年卡" {
		t.Fatalf("MemberCardTypeText(year) = %s, want 年卡", got)
	}
	if got := constants.MemberCardMonths(constants.MemberCardTypeYear); got != 12 {
		t.Fatalf("MemberCardMonths(year) = %d, want 12", got)
	}
	if got := constants.MemberCardPrice(constants.MemberCardTypeYear); got != 899 {
		t.Fatalf("MemberCardPrice(year) = %f, want 899", got)
	}
}
