package handler

import (
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

// ActivityHandler exposes activity endpoints.
type ActivityHandler struct {
	svc    *service.ActivityService
	users  *service.UserService
	logger *slog.Logger
}

// NewActivityHandler wires the activity handler dependencies.
func NewActivityHandler(svc *service.ActivityService, users *service.UserService, logger *slog.Logger) *ActivityHandler {
	return &ActivityHandler{svc: svc, users: users, logger: logger}
}

// Create handles POST /activities (admin).
func (h *ActivityHandler) Create(c *gin.Context) {
	var req dto.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	a, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

// List handles GET /activities.
func (h *ActivityHandler) List(c *gin.Context) {
	var q dto.PageQuery
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

// Get handles GET /activities/:id.
func (h *ActivityHandler) Get(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动ID不合法")
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

// Signup handles POST /activities/:id/signup.
func (h *ActivityHandler) Signup(c *gin.Context) {
	user := h.requireUser(c)
	if user == nil {
		return
	}
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "活动ID不合法")
		return
	}
	reg, err := h.svc.Signup(c.Request.Context(), user, id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, reg)
}

func (h *ActivityHandler) requireUser(c *gin.Context) *model.User {
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
