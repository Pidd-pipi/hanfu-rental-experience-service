package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/handler"
)

// RegisterActivityRoutes registers activity endpoints.
func RegisterActivityRoutes(g *gin.RouterGroup, h *handler.ActivityHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	activities := g.Group("/activities")
	{
		activities.GET("", apiLimiter, h.List)
		activities.GET("/:id", apiLimiter, h.Get)
		authed := activities.Group("", auth)
		{
			authed.POST("/:id/signup", apiLimiter, h.Signup)
		}
		admin := activities.Group("", auth, requireAdmin)
		{
			admin.POST("", apiLimiter, h.Create)
		}
	}
}
