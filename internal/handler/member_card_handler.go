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

// MemberCardHandler exposes member card endpoints.
type MemberCardHandler struct {
	svc    *service.MemberCardService
	logger *slog.Logger
}

// NewMemberCardHandler wires the member card handler dependencies.
func NewMemberCardHandler(svc *service.MemberCardService, logger *slog.Logger) *MemberCardHandler {
	return &MemberCardHandler{svc: svc, logger: logger}
}

// Create handles POST /member-cards.
func (h *MemberCardHandler) Create(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	var req dto.CreateMemberCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	card, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, card)
}

// ListMy handles GET /member-cards/me.
func (h *MemberCardHandler) ListMy(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	items, err := h.svc.ListMy(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}
