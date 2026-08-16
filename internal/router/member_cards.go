package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/handler"
)

// RegisterMemberCardRoutes registers member card endpoints.
func RegisterMemberCardRoutes(g *gin.RouterGroup, h *handler.MemberCardHandler, auth, apiLimiter gin.HandlerFunc) {
	cards := g.Group("/member-cards", auth)
	{
		cards.POST("", apiLimiter, h.Create)
		cards.GET("/me", apiLimiter, h.ListMy)
	}
}
