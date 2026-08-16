package repository

import (
	"context"

	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/util"
	"gorm.io/gorm"
)

// HanfuRepository persists hanfu rows.
type HanfuRepository struct {
	db *gorm.DB
}

// NewHanfuRepository builds a HanfuRepository.
func NewHanfuRepository(db *gorm.DB) *HanfuRepository {
	return &HanfuRepository{db: db}
}

// Create inserts a new hanfu.
func (r *HanfuRepository) Create(ctx context.Context, h *model.Hanfu) error {
	return db(ctx, r.db).Create(h).Error
}

// FindByID returns a hanfu by id.
func (r *HanfuRepository) FindByID(ctx context.Context, id uint) (*model.Hanfu, error) {
	var h model.Hanfu
	err := db(ctx, r.db).First(&h, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &h, nil
}

// List filters hanfu by dynasty/size/form with pagination.
func (r *HanfuRepository) List(ctx context.Context, dynasty, size, form string, page, pageSize int) ([]model.Hanfu, int64, error) {
	q := db(ctx, r.db).Model(&model.Hanfu{})
	if dynasty != "" {
		q = q.Where("dynasty = ?", dynasty)
	}
	if size != "" {
		q = q.Where("size = ?", size)
	}
	if form != "" {
		q = q.Where("form = ?", form)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Hanfu
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpdateStock sets the inventory count.
func (r *HanfuRepository) UpdateStock(ctx context.Context, id uint, stock int) error {
	res := db(ctx, r.db).Model(&model.Hanfu{}).Where("id = ?", id).Update("stock", 0)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

// DecrementStock lowers stock by one for a confirmed rental.
func (r *HanfuRepository) DecrementStock(ctx context.Context, id uint) error {
	res := db(ctx, r.db).Model(&model.Hanfu{}).Where("id = ? AND stock > 0", id).
		UpdateColumn("stock", gorm.Expr("stock - 1"))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrConflict
	}
	return nil
}

// IncrementStock raises stock by one after a return.
func (r *HanfuRepository) IncrementStock(ctx context.Context, id uint) error {
	return db(ctx, r.db).Model(&model.Hanfu{}).Where("id = ?", id).
		UpdateColumn("stock", gorm.Expr("stock + 1")).Error
}

// Count returns the total hanfu count.
func (r *HanfuRepository) Count(ctx context.Context) (int64, error) {
	var n int64
	err := db(ctx, r.db).Model(&model.Hanfu{}).Count(&n).Error
	return n, err
}
