package repository

import (
	"context"

	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/util"
	"gorm.io/gorm"
)

// RentalOrderRepository persists rental order rows.
type RentalOrderRepository struct {
	db *gorm.DB
}

// NewRentalOrderRepository builds a RentalOrderRepository.
func NewRentalOrderRepository(db *gorm.DB) *RentalOrderRepository {
	return &RentalOrderRepository{db: db}
}

// Transaction runs fn inside a database transaction for cross-repository writes.
func (r *RentalOrderRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return Transaction(ctx, r.db, fn)
}

// Create inserts a new rental order.
func (r *RentalOrderRepository) Create(ctx context.Context, o *model.RentalOrder) error {
	return db(ctx, r.db).Create(o).Error
}

// FindByID returns a rental order by id.
func (r *RentalOrderRepository) FindByID(ctx context.Context, id uint) (*model.RentalOrder, error) {
	var o model.RentalOrder
	err := db(ctx, r.db).First(&o, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &o, nil
}

// ListByUser returns the rental orders of a user.
func (r *RentalOrderRepository) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]model.RentalOrder, int64, error) {
	q := db(ctx, r.db).Model(&model.RentalOrder{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.RentalOrder
	err := q.Order("created_at ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ListAll returns all rental orders (admin analytics).
func (r *RentalOrderRepository) ListAll(ctx context.Context, page, pageSize int) ([]model.RentalOrder, int64, error) {
	q := db(ctx, r.db).Model(&model.RentalOrder{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.RentalOrder
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpdateStatus sets the order status.
func (r *RentalOrderRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	res := db(ctx, r.db).Model(&model.RentalOrder{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

// UpdateStatusIf performs an atomic order status transition.
func (r *RentalOrderRepository) UpdateStatusIf(ctx context.Context, id uint, from, to string) error {
	res := db(ctx, r.db).Model(&model.RentalOrder{}).Where("id = ? AND status = ?", id, from).Update("status", to)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrConflict
	}
	return nil
}
