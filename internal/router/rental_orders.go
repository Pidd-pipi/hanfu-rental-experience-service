package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/handler"
)

// RegisterRentalOrderRoutes registers rental order endpoints.
func RegisterRentalOrderRoutes(g *gin.RouterGroup, h *handler.RentalOrderHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	orders := g.Group("/rental-orders")
	{
		me := orders.Group("/me", auth)
		{
			me.GET("", apiLimiter, h.ListMy)
		}
		authed := orders.Group("", auth)
		{
			authed.POST("", apiLimiter, h.Create)
			authed.POST("/:id/cancel", apiLimiter, h.Cancel)
		}
		admin := orders.Group("", auth, requireAdmin)
		{
			admin.GET("", apiLimiter, h.ListAll)
			admin.POST("/:id/confirm", apiLimiter, h.Confirm)
			admin.POST("/:id/return", apiLimiter, h.Return)
			admin.POST("/:id/complete", apiLimiter, h.Complete)
		}
	}
}
