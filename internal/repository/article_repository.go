package repository

import (
	"context"

	"github.com/lp/hanfu-rental/internal/model"
	"gorm.io/gorm"
)

// ArticleRepository persists article rows.
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository builds an ArticleRepository.
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create inserts a new article.
func (r *ArticleRepository) Create(ctx context.Context, a *model.Article) error {
	return db(ctx, r.db).Create(a).Error
}

// FindByID returns an article by id.
func (r *ArticleRepository) FindByID(ctx context.Context, id uint) (*model.Article, error) {
	var a model.Article
	err := db(ctx, r.db).First(&a, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &a, nil
}

// List returns articles with pagination.
func (r *ArticleRepository) List(ctx context.Context, page, pageSize int) ([]model.Article, int64, error) {
	q := db(ctx, r.db).Model(&model.Article{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Article
	err := q.Order("published_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
