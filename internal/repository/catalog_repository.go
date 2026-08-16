package repository

import (
	"context"

	"github.com/lp/hanfu-rental/internal/model"
	"gorm.io/gorm"
)

// CatalogRepository persists makeup packages and photographers.
type CatalogRepository struct {
	db *gorm.DB
}

// NewCatalogRepository builds a CatalogRepository.
func NewCatalogRepository(db *gorm.DB) *CatalogRepository {
	return &CatalogRepository{db: db}
}

// ListMakeupPackages returns all makeup packages.
func (r *CatalogRepository) ListMakeupPackages(ctx context.Context) ([]model.MakeupPackage, error) {
	var items []model.MakeupPackage
	err := db(ctx, r.db).Order("price ASC").Find(&items).Error
	return items, err
}

// FindMakeupPackage returns a makeup package by id.
func (r *CatalogRepository) FindMakeupPackage(ctx context.Context, id uint) (*model.MakeupPackage, error) {
	var p model.MakeupPackage
	err := db(ctx, r.db).First(&p, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &p, nil
}

// ListPhotographers returns all photographers.
func (r *CatalogRepository) ListPhotographers(ctx context.Context) ([]model.Photographer, error) {
	var items []model.Photographer
	err := db(ctx, r.db).Order("price_per_day ASC").Find(&items).Error
	return items, err
}

// FindPhotographer returns a photographer by id.
func (r *CatalogRepository) FindPhotographer(ctx context.Context, id uint) (*model.Photographer, error) {
	var p model.Photographer
	err := db(ctx, r.db).First(&p, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &p, nil
}
