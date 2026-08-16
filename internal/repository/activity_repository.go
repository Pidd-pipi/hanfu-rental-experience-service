package repository

import (
	"context"

	"github.com/lp/hanfu-rental/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ActivityRepository persists activity and registration rows.
type ActivityRepository struct {
	db *gorm.DB
}

// NewActivityRepository builds an ActivityRepository.
func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

// Transaction runs fn inside a database transaction for cross-repository writes.
func (r *ActivityRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return Transaction(ctx, r.db, fn)
}

// Create inserts a new activity.
func (r *ActivityRepository) Create(ctx context.Context, a *model.Activity) error {
	return db(ctx, r.db).Create(a).Error
}

// FindByID returns an activity by id.
func (r *ActivityRepository) FindByID(ctx context.Context, id uint) (*model.Activity, error) {
	var a model.Activity
	err := db(ctx, r.db).First(&a, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &a, nil
}

// FindByIDForUpdate locks the activity row (SELECT ... FOR UPDATE) to serialize
// signup capacity checks inside a transaction.
func (r *ActivityRepository) FindByIDForUpdate(ctx context.Context, id uint) (*model.Activity, error) {
	var a model.Activity
	err := db(ctx, r.db).Clauses(clause.Locking{Strength: "UPDATE"}).First(&a, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &a, nil
}

// List returns activities with pagination.
func (r *ActivityRepository) List(ctx context.Context, page, pageSize int) ([]model.Activity, int64, error) {
	q := db(ctx, r.db).Model(&model.Activity{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Activity
	err := q.Order("start_time ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// CountRegistrations returns the number of signups for an activity.
func (r *ActivityRepository) CountRegistrations(ctx context.Context, activityID uint) (int64, error) {
	var n int64
	err := db(ctx, r.db).Model(&model.ActivityRegistration{}).Where("activity_id = ?", activityID).Count(&n).Error
	return n, err
}

// FindRegistration returns a signup of a user for an activity.
func (r *ActivityRepository) FindRegistration(ctx context.Context, activityID, userID uint) (*model.ActivityRegistration, error) {
	var reg model.ActivityRegistration
	err := db(ctx, r.db).Where("activity_id = ? AND user_id = ?", activityID, userID).First(&reg).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &reg, nil
}

// CreateRegistration inserts a signup.
func (r *ActivityRepository) CreateRegistration(ctx context.Context, reg *model.ActivityRegistration) error {
	return db(ctx, r.db).Create(reg).Error
}
