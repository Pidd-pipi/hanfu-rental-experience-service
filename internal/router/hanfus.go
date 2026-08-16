package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/handler"
)

// RegisterHanfuRoutes registers hanfu catalog endpoints.
func RegisterHanfuRoutes(g *gin.RouterGroup, h *handler.HanfuHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	hanfus := g.Group("/hanfus")
	{
		hanfus.GET("", apiLimiter, h.List)
		hanfus.GET("/:id", apiLimiter, h.Get)
		admin := hanfus.Group("", auth, requireAdmin)
		{
			admin.POST("", apiLimiter, h.Create)
			admin.PUT("/:id/stock", apiLimiter, h.UpdateStock)
		}
	}
}
