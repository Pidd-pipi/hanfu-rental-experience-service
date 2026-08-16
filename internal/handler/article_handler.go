package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/dto"
	"github.com/lp/hanfu-rental/internal/service"
	"github.com/lp/hanfu-rental/internal/util"
)

// ArticleHandler exposes article endpoints.
type ArticleHandler struct {
	svc    *service.ArticleService
	logger *slog.Logger
}

// NewArticleHandler wires the article handler dependencies.
func NewArticleHandler(svc *service.ArticleService, logger *slog.Logger) *ArticleHandler {
	return &ArticleHandler{svc: svc, logger: logger}
}

// Create handles POST /articles (admin).
func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
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

// List handles GET /articles.
func (h *ArticleHandler) List(c *gin.Context) {
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

// Get handles GET /articles/:id.
func (h *ArticleHandler) Get(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章ID不合法")
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}
