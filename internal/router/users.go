package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/handler"
)

// RegisterUserRoutes registers customer endpoints.
func RegisterUserRoutes(g *gin.RouterGroup, h *handler.UserHandler, auth, requireAdmin, loginLimiter, apiLimiter gin.HandlerFunc) {
	users := g.Group("/users")
	{
		users.POST("/register", loginLimiter, h.Register)
		users.POST("/login", loginLimiter, h.Login)
		me := users.Group("/me", auth)
		{
			me.GET("", apiLimiter, h.GetProfile)
			me.PUT("", apiLimiter, h.UpdateProfile)
			me.POST("/deposit", apiLimiter, h.PayDeposit)
		}
		admin := users.Group("", auth, requireAdmin)
		admin.POST("/:id/deposit/refund", apiLimiter, h.RefundDeposit)
	}
}
