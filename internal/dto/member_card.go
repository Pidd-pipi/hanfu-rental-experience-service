package dto

// CreateMemberCardRequest is the payload for applying a member card.
type CreateMemberCardRequest struct {
	CardType string `json:"card_type" binding:"required,oneof=month quarter year"`
}
