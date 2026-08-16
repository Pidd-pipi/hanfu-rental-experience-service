package util

import (
	"fmt"
	"time"

	"github.com/lp/hanfu-rental/internal/constants"
)

// FormatDateTime renders a time value using the shared display layout.
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// FormatDate renders a date-only display string.
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatMoney renders a money value with two decimals.
func FormatMoney(v float64) string {
	return fmt.Sprintf("¥%.2f", v)
}

// RentalStatusText maps a rental status to its Chinese label.
func RentalStatusText(s string) string {
	return constants.RentalOrderStatusText(s)
}

// RoleText maps a user role to its Chinese label.
func RoleText(role string) string {
	return constants.UserRoleText(role)
}

// DynastyText maps a hanfu dynasty to its Chinese label.
func DynastyText(d string) string {
	return constants.HanfuDynastyText(d)
}

// CardTypeText maps a member card type to its Chinese label.
func CardTypeText(t string) string {
	return constants.MemberCardTypeText(t)
}
