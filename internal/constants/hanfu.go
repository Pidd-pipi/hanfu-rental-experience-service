// Package constants centralizes business constants for the hanfu-rental backend.
package constants

// HanfuDynasty defines the hanfu dynasty enum values shared with the frontend.
const (
	HanfuDynastyHan  = "han"
	HanfuDynastyTang = "tang"
	HanfuDynastySong = "song"
	HanfuDynastyMing = "ming"
)

// HanfuDynasties lists all valid dynasties.
var HanfuDynasties = []string{HanfuDynastyHan, HanfuDynastyTang, HanfuDynastySong, HanfuDynastyMing}

// IsHanfuDynasty reports whether the given dynasty is valid.
func IsHanfuDynasty(d string) bool {
	for _, v := range HanfuDynasties {
		if v == d {
			return true
		}
	}
	return false
}

// HanfuStatus defines hanfu availability states.
const (
	HanfuStatusAvailable = "available"
	HanfuStatusRented    = "rented"
	HanfuStatusMaintenance = "maintenance"
)

// HanfuDynastyText returns the Chinese label of a dynasty.
func HanfuDynastyText(d string) string {
	switch d {
	case HanfuDynastyHan:
		return "汉制"
	case HanfuDynastyTang:
		return "唐制"
	case HanfuDynastySong:
		return "宋制"
	case HanfuDynastyMing:
		return "明制"
	default:
		return "未知"
	}
}
