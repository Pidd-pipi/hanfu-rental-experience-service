package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/repository"
	"github.com/lp/hanfu-rental/internal/util"
)

// CatalogHandler exposes makeup packages and photographers.
type CatalogHandler struct {
	catalog *repository.CatalogRepository
	logger  *slog.Logger
}

// NewCatalogHandler wires the catalog handler dependencies.
func NewCatalogHandler(catalog *repository.CatalogRepository, logger *slog.Logger) *CatalogHandler {
	return &CatalogHandler{catalog: catalog, logger: logger}
}

// ListMakeupPackages handles GET /makeup-packages.
func (h *CatalogHandler) ListMakeupPackages(c *gin.Context) {
	items, err := h.catalog.ListMakeupPackages(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}

// ListPhotographers handles GET /photographers.
func (h *CatalogHandler) ListPhotographers(c *gin.Context) {
	items, err := h.catalog.ListPhotographers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}

var _ = constants.CodeInternalError
