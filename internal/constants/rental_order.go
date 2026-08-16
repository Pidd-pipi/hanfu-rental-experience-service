package constants

// RentalOrderStatus defines rental order state machine values shared with the frontend.
const (
	RentalOrderStatusPending   = "pending"
	RentalOrderStatusRenting   = "renting"
	RentalOrderStatusReturned  = "returned"
	RentalOrderStatusCompleted = "completed"
	RentalOrderStatusCancelled = "cancelled"
)

// RentalOrderStatuses lists all valid rental statuses in flow order.
var RentalOrderStatuses = []string{
	RentalOrderStatusPending, RentalOrderStatusRenting, RentalOrderStatusReturned,
	RentalOrderStatusCompleted, RentalOrderStatusCancelled,
}

// IsRentalOrderStatus reports whether the given status is valid.
func IsRentalOrderStatus(s string) bool {
	for _, v := range RentalOrderStatuses {
		if v == s {
			return true
		}
	}
	return false
}

// RentalOrderStatusText returns the Chinese label of a rental status.
func RentalOrderStatusText(s string) string {
	switch s {
	case RentalOrderStatusPending:
		return "待确认"
	case RentalOrderStatusRenting:
		return "租赁中"
	case RentalOrderStatusReturned:
		return "已归还"
	case RentalOrderStatusCompleted:
		return "已完成"
	case RentalOrderStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}
