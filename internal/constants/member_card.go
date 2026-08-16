package constants

// MemberCardType defines member card type enum values shared with the frontend.
const (
	MemberCardTypeMonth  = "month"
	MemberCardTypeQuarter = "quarter"
	MemberCardTypeYear   = "year"
)

// MemberCardTypes lists all valid card types.
var MemberCardTypes = []string{MemberCardTypeMonth, MemberCardTypeQuarter, MemberCardTypeYear}

// IsMemberCardType reports whether the given card type is valid.
func IsMemberCardType(t string) bool {
	for _, v := range MemberCardTypes {
		if v == t {
			return true
		}
	}
	return false
}

// MemberCardStatus defines card lifecycle states.
const (
	MemberCardStatusActive   = "active"
	MemberCardStatusExpired  = "expired"
	MemberCardStatusDisabled = "disabled"
)

// MemberCardTypeText returns the Chinese label of a card type.
func MemberCardTypeText(t string) string {
	switch t {
	case MemberCardTypeMonth:
		return "月卡"
	case MemberCardTypeQuarter:
		return "季卡"
	case MemberCardTypeYear:
		return "年卡"
	default:
		return "未知"
	}
}

// MemberCardMonths maps a card type to its validity months.
func MemberCardMonths(t string) int {
	switch t {
	case MemberCardTypeMonth:
		return 1
	case MemberCardTypeQuarter:
		return 3
	case MemberCardTypeYear:
		return 12
	default:
		return 0
	}
}

// MemberCardPrice maps a card type to its price.
func MemberCardPrice(t string) float64 {
	switch t {
	case MemberCardTypeMonth:
		return 99
	case MemberCardTypeQuarter:
		return 269
	case MemberCardTypeYear:
		return 899
	default:
		return 0
	}
}
