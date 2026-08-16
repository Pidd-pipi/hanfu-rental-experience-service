package dto

// RegisterRequest is the customer registration payload.
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"required,min=2,max=32"`
}

// LoginRequest is the login payload.
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

// LoginResponse carries the JWT and the basic profile.
type LoginResponse struct {
	Token string    `json:"token"`
	User  *UserView `json:"user"`
}

// UserView is the public user profile view.
type UserView struct {
	ID       uint    `json:"id"`
	Phone    string  `json:"phone"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Role     string  `json:"role"`
	Deposit  float64 `json:"deposit"`
}

// UpdateProfileRequest is the profile update payload.
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=32"`
	Avatar   string `json:"avatar"`
}

// PayDepositRequest is the deposit payment payload.
type PayDepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
