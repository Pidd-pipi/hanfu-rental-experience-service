package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/service"
	"github.com/lp/hanfu-rental/internal/util"
)

// HanfuHandler exposes hanfu catalog endpoints.
type HanfuHandler struct {
	svc    *service.HanfuService
	logger *slog.Logger
}

// NewHanfuHandler wires the hanfu handler dependencies.
func NewHanfuHandler(svc *service.HanfuService, logger *slog.Logger) *HanfuHandler {
	return &HanfuHandler{svc: svc, logger: logger}
}

// Create handles POST /hanfus (admin).
func (h *HanfuHandler) Create(c *gin.Context) {
	var req dto.CreateHanfuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	hanfu, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(fmt.Errorf("handler hanfu create: %v", err))
		return
	}
	util.OK(c, hanfu)
}

// Get handles GET /hanfus/:id.
func (h *HanfuHandler) Get(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "汉服ID不合法")
		return
	}
	hanfu, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, hanfu)
}

// List handles GET /hanfus.
func (h *HanfuHandler) List(c *gin.Context) {
	var q dto.ListHanfuQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	result, err := h.svc.List(c.Request.Context(), &q)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// UpdateStock handles PUT /hanfus/:id/stock (admin).
func (h *HanfuHandler) UpdateStock(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "汉服ID不合法")
		return
	}
	var req struct {
		Stock int `json:"stock" binding:"gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	hanfu, err := h.svc.UpdateStock(c.Request.Context(), id, req.Stock)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, hanfu)
}
