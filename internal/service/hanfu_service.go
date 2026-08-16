package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/util"
)

// HanfuRepository is the data access contract for hanfu rows.
type HanfuRepository interface {
	Create(ctx context.Context, h *model.Hanfu) error
	FindByID(ctx context.Context, id uint) (*model.Hanfu, error)
	List(ctx context.Context, dynasty, size, form string, page, pageSize int) ([]model.Hanfu, int64, error)
	UpdateStock(ctx context.Context, id uint, stock int) error
	DecrementStock(ctx context.Context, id uint) error
	IncrementStock(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}

// HanfuService manages the hanfu rental catalog.
type HanfuService struct {
	hanfus HanfuRepository
	logger *slog.Logger
}

// NewHanfuService wires the hanfu service dependencies.
func NewHanfuService(hanfus HanfuRepository, logger *slog.Logger) *HanfuService {
	return &HanfuService{hanfus: hanfus, logger: logger}
}

// Create adds a new hanfu.
func (s *HanfuService) Create(ctx context.Context, req *dto.CreateHanfuRequest) (*model.Hanfu, error) {
	h := &model.Hanfu{
		Name: req.Name, Dynasty: req.Dynasty, Form: req.Form, Color: req.Color,
		Size: req.Size, PricePerDay: req.PricePerDay, Images: req.Images,
		Stock: req.Stock, Status: constants.HanfuStatusAvailable,
	}
	if err := s.hanfus.Create(ctx, h); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogHanfuCreateFailed, req.Name, err))
		return nil, util.WrapAppError(fmt.Errorf("hanfu[name=%s] create: %w", req.Name, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogHanfuCreateSuccess, h.ID, h.Name, h.Dynasty))
	return h, nil
}

// Get returns one hanfu.
func (s *HanfuService) Get(ctx context.Context, id uint) (*model.Hanfu, error) {
	h, err := s.hanfus.FindByID(ctx, id)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("hanfu[id=%d] get: %w", id, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	return h, nil
}

// List filters hanfu.
func (s *HanfuService) List(ctx context.Context, q *dto.ListHanfuQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.hanfus.List(ctx, q.Dynasty, "", "", q.Page, 0)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("hanfu list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// UpdateStock maintains inventory.
func (s *HanfuService) UpdateStock(ctx context.Context, id uint, stock int) (*model.Hanfu, error) {
	if err := s.hanfus.UpdateStock(ctx, id, stock); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("hanfu[id=%d] stock update: %w", id, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	s.logger.Info(fmt.Sprintf(constants.LogHanfuStockUpdateSuccess, id, stock))
	return s.hanfus.FindByID(ctx, id)
}
