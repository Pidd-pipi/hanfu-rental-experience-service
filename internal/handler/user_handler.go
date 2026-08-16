// Package handler implements the HTTP handlers for hanfu-rental.
package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/middleware"
	"github.com/lp/hanfu-rental/internal/service"
	"github.com/lp/hanfu-rental/internal/util"
)

// UserHandler exposes customer endpoints.
type UserHandler struct {
	svc    *service.UserService
	logger *slog.Logger
}

// NewUserHandler wires the user handler dependencies.
func NewUserHandler(svc *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

// Register handles POST /users/register.
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	user, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, service.ToUserView(user))
}

// Login handles POST /users/login.
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, resp)
}

// GetProfile handles GET /users/me.
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	_, err = h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, nil)
}

// UpdateProfile handles PUT /users/me.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, service.ToUserView(user))
}

// PayDeposit handles POST /users/me/deposit.
func (h *UserHandler) PayDeposit(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	var req dto.PayDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	user, err := h.svc.PayDeposit(c.Request.Context(), userID, req.Amount)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, service.ToUserView(user))
}

// RefundDeposit handles POST /users/:id/deposit/refund (admin).
func (h *UserHandler) RefundDeposit(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户ID不合法")
		return
	}
	var req dto.PayDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	user, err := h.svc.RefundDeposit(c.Request.Context(), id, req.Amount)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, service.ToUserView(user))
}
