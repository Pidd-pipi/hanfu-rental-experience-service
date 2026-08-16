package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/middleware"
	"github.com/lp/hanfu-rental/internal/model"
	"github.com/lp/hanfu-rental/internal/service"
	"github.com/lp/hanfu-rental/internal/util"
)

// RentalOrderHandler exposes rental order endpoints.
type RentalOrderHandler struct {
	svc    *service.RentalOrderService
	users  *service.UserService
	logger *slog.Logger
}

// NewRentalOrderHandler wires the rental order handler dependencies.
func NewRentalOrderHandler(svc *service.RentalOrderService, users *service.UserService, logger *slog.Logger) *RentalOrderHandler {
	return &RentalOrderHandler{svc: svc, users: users, logger: logger}
}

// Create handles POST /rental-orders.
func (h *RentalOrderHandler) Create(c *gin.Context) {
	user := h.requireUser(c)
	if user == nil {
		return
	}
	var req dto.CreateRentalOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	order, err := h.svc.Create(c.Request.Context(), user, &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, order)
}

// ListMy handles GET /rental-orders/me.
func (h *RentalOrderHandler) ListMy(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	var q dto.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	result, err := h.svc.ListMy(c.Request.Context(), userID, &q)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// ListAll handles GET /rental-orders (admin).
func (h *RentalOrderHandler) ListAll(c *gin.Context) {
	var q dto.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	result, err := h.svc.ListAll(c.Request.Context(), &q)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Confirm handles POST /rental-orders/:id/confirm (admin).
func (h *RentalOrderHandler) Confirm(c *gin.Context) {
	h.act(c, h.svc.Confirm)
}

// Return handles POST /rental-orders/:id/return (admin).
func (h *RentalOrderHandler) Return(c *gin.Context) {
	h.act(c, h.svc.Return)
}

// Complete handles POST /rental-orders/:id/complete (admin).
func (h *RentalOrderHandler) Complete(c *gin.Context) {
	h.act(c, h.svc.Complete)
}

// Cancel handles POST /rental-orders/:id/cancel.
func (h *RentalOrderHandler) Cancel(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "订单ID不合法")
		return
	}
	order, err := h.svc.Cancel(c.Request.Context(), userID, id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, order)
}

func (h *RentalOrderHandler) act(c *gin.Context, fn func(ctx context.Context, orderID uint) (*model.RentalOrder, error)) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "订单ID不合法")
		return
	}
	order, err := fn(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, order)
}

func (h *RentalOrderHandler) requireUser(c *gin.Context) *model.User {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return nil
	}
	user, err := h.users.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return nil
	}
	return user
}
