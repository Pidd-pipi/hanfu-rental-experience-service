package constants

// UserRole defines user role enum values shared with the frontend.
const (
	UserRoleCustomer = "customer"
	UserRoleAdmin    = "admin"
)

// UserRoles lists all valid user roles.
var UserRoles = []string{UserRoleCustomer, UserRoleAdmin}

// IsUserRole reports whether the given role is valid.
func IsUserRole(r string) bool {
	for _, v := range UserRoles {
		if v == r {
			return true
		}
	}
	return false
}

// UserRoleText returns the Chinese label of a user role.
func UserRoleText(r string) string {
	switch r {
	case UserRoleCustomer:
		return "顾客"
	case UserRoleAdmin:
		return "管理员"
	default:
		return "未知"
	}
}
