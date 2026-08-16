package repository

import (
	"context"

	"github.com/lp/hanfu-rental/internal/model"
	"gorm.io/gorm"
)

// MemberCardRepository persists member card rows.
type MemberCardRepository struct {
	db *gorm.DB
}

// NewMemberCardRepository builds a MemberCardRepository.
func NewMemberCardRepository(db *gorm.DB) *MemberCardRepository {
	return &MemberCardRepository{db: db}
}

// Create inserts a new member card.
func (r *MemberCardRepository) Create(ctx context.Context, c *model.MemberCard) error {
	return db(ctx, r.db).Create(c).Error
}

// FindActiveByUser returns the active card of a user.
func (r *MemberCardRepository) FindActiveByUser(ctx context.Context, userID uint) (*model.MemberCard, error) {
	var c model.MemberCard
	err := db(ctx, r.db).Where("user_id = ? AND status = ?", userID, "active").Order("created_at DESC").First(&c).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &c, nil
}

// ListByUser returns all cards of a user.
func (r *MemberCardRepository) ListByUser(ctx context.Context, userID uint) ([]model.MemberCard, error) {
	var items []model.MemberCard
	err := db(ctx, r.db).Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, err
}

// Expire marks cards whose expire_at passed as expired.
func (r *MemberCardRepository) Expire(ctx context.Context, now interface{}) error {
	return db(ctx, r.db).Model(&model.MemberCard{}).
		Where("status = ? AND expire_at < ?", "active", now).Update("status", "expired").Error
}
