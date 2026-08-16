package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/repository"
	"github.com/lp/hanfu-rental/internal/util"
)

// RentalOrderService manages combined hanfu rental orders.
type RentalOrderService struct {
	orders  *repository.RentalOrderRepository
	hanfus  *repository.HanfuRepository
	catalog *repository.CatalogRepository
	cards   *repository.MemberCardRepository
	users   *repository.UserRepository
	logger  *slog.Logger
}

// NewRentalOrderService wires the rental order service dependencies.
func NewRentalOrderService(orders *repository.RentalOrderRepository, hanfus *repository.HanfuRepository, catalog *repository.CatalogRepository, cards *repository.MemberCardRepository, users *repository.UserRepository, logger *slog.Logger) *RentalOrderService {
	return &RentalOrderService{orders: orders, hanfus: hanfus, catalog: catalog, cards: cards, users: users, logger: logger}
}

// Create builds a pending combined order (hanfu + makeup + photographer).
func (s *RentalOrderService) Create(ctx context.Context, user *model.User, req *dto.CreateRentalOrderRequest) (*model.RentalOrder, error) {
	if !req.EndDate.After(req.StartDate) {
		return nil, util.NewAppError(400, constants.CodeBadRequest, "归还日期必须晚于开始日期", nil)
	}
	hanfu, err := s.hanfus.FindByID(ctx, req.HanfuID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order[user=%d] hanfu lookup: %v", user.ID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if hanfu.Stock <= 0 || hanfu.Status != constants.HanfuStatusAvailable {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgHanfuUnavailable, nil)
	}
	makeupFee := 0.0
	if req.MakeupPackageID != nil {
		pkg, err := s.catalog.FindMakeupPackage(ctx, *req.MakeupPackageID)
		if err != nil {
			return nil, util.WrapAppError(fmt.Errorf("rental_order[user=%d] makeup lookup: %w", user.ID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
		}
		makeupFee = pkg.Price
	}
	photographerFee := 0.0
	if req.PhotographerID != nil {
		p, err := s.catalog.FindPhotographer(ctx, *req.PhotographerID)
		if err != nil {
			return nil, util.WrapAppError(fmt.Errorf("rental_order[user=%d] photographer lookup: %w", user.ID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
		}
		photographerFee = p.PricePerDay * float64(util.CalcRentalDays(req.StartDate, req.EndDate))
	}
	hasCard := false
	cardType := ""
	if card, err := s.cards.FindActiveByUser(ctx, user.ID); err == nil && card != nil {
		hasCard = true
		cardType = card.CardType
	}
	fee := util.CalcRentalFee(hanfu.PricePerDay, req.StartDate, req.EndDate, hasCard, cardType, makeupFee, photographerFee)
	s.logger.Info(fmt.Sprintf(constants.LogFeeCalcUsed, fee.Days, fee.Subtotal, fee.CardDiscount))
	order := &model.RentalOrder{
		UserID: user.ID, HanfuID: req.HanfuID, MakeupPackageID: req.MakeupPackageID,
		PhotographerID: req.PhotographerID, StartDate: req.StartDate, EndDate: req.EndDate,
		DepositAmount: util.CalcDeposit(hanfu.PricePerDay), RentalFee: fee.Total,
		Status: constants.RentalOrderStatusPending,
	}
	if err := s.orders.Create(ctx, order); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogRentalOrderCreateFailed, user.ID, req.HanfuID, err))
		return nil, util.WrapAppError(fmt.Errorf("rental_order[user=%d] create: %w", user.ID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogRentalOrderCreateSuccess, order.ID, user.ID, req.HanfuID))
	return order, nil
}

// ListMy returns the rental orders of the current user.
func (s *RentalOrderService) ListMy(ctx context.Context, userID uint, q *dto.PageQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.orders.ListByUser(ctx, userID, q.Page, q.PageSize)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order[user=%d] list: %w", userID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// ListAll returns all orders (admin).
func (s *RentalOrderService) ListAll(ctx context.Context, q *dto.PageQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.orders.ListAll(ctx, q.Page, q.PageSize)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order list all: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// Confirm approves a pending order (admin) and deducts hanfu stock.
func (s *RentalOrderService) Confirm(ctx context.Context, orderID uint) (*model.RentalOrder, error) {
	order, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] confirm find: %w", orderID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if order.Status != constants.RentalOrderStatusPending {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgOrderStatusInvalid, nil)
	}
	if err := s.orders.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.orders.UpdateStatusIf(txCtx, orderID, constants.RentalOrderStatusPending, constants.RentalOrderStatusRenting); err != nil {
			return err
		}
		if err := s.hanfus.DecrementStock(txCtx, order.HanfuID); err != nil {
			if errors.Is(err, util.ErrConflict) {
				return util.NewAppError(409, constants.CodeConflict, constants.MsgHanfuUnavailable, err)
			}
			return err
		}
		return nil
	}); err != nil {
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgOrderStatusInvalid, nil)
		}
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] confirm: %w", orderID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogRentalOrderConfirmSuccess, orderID))
	order.Status = constants.RentalOrderStatusRenting
	return order, nil
}

// Return marks a renting order as returned (admin) and restores stock.
func (s *RentalOrderService) Return(ctx context.Context, orderID uint) (*model.RentalOrder, error) {
	order, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] return find: %w", orderID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if order.Status != constants.RentalOrderStatusRenting {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgOrderStatusInvalid, nil)
	}
	if err := s.orders.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.orders.UpdateStatusIf(txCtx, orderID, constants.RentalOrderStatusRenting, constants.RentalOrderStatusReturned); err != nil {
			return err
		}
		if err := s.hanfus.IncrementStock(txCtx, order.HanfuID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgOrderStatusInvalid, nil)
		}
		s.logger.Error(fmt.Sprintf(constants.LogRentalOrderReturnFailed, orderID, err))
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] return: %w", orderID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogRentalOrderReturnSuccess, orderID))
	order.Status = constants.RentalOrderStatusReturned
	return order, nil
}

// Complete finishes a returned order and refunds the deposit.
func (s *RentalOrderService) Complete(ctx context.Context, orderID uint) (*model.RentalOrder, error) {
	order, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] complete find: %w", orderID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if order.Status != constants.RentalOrderStatusReturned {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgOrderStatusInvalid, nil)
	}
	if err := s.orders.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.orders.UpdateStatusIf(txCtx, orderID, constants.RentalOrderStatusReturned, constants.RentalOrderStatusCompleted); err != nil {
			return err
		}
		if err := s.users.AddDeposit(txCtx, order.UserID, order.DepositAmount); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgOrderStatusInvalid, nil)
		}
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] complete: %w", orderID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogRentalOrderCompleteSuccess, orderID, order.DepositAmount))
	order.Status = constants.RentalOrderStatusCompleted
	return order, nil
}

// Cancel cancels a pending order.
func (s *RentalOrderService) Cancel(ctx context.Context, userID, orderID uint) (*model.RentalOrder, error) {
	order, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] cancel find: %w", orderID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if order.UserID != userID {
		return nil, util.NewAppError(403, constants.CodeForbidden, constants.MsgNotOwner, nil)
	}
	if order.Status != constants.RentalOrderStatusPending {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgOrderStatusInvalid, nil)
	}
	if err := s.orders.UpdateStatus(ctx, orderID, constants.RentalOrderStatusCancelled); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("rental_order[id=%d] cancel: %w", orderID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogRentalOrderCancelSuccess, orderID))
	order.Status = constants.RentalOrderStatusCancelled
	return order, nil
}
