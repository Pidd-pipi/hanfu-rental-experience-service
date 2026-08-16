package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/handler"
)

// RegisterArticleRoutes registers article endpoints.
func RegisterArticleRoutes(g *gin.RouterGroup, h *handler.ArticleHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	articles := g.Group("/articles")
	{
		articles.GET("", apiLimiter, h.List)
		articles.GET("/:id", apiLimiter, h.Get)
		admin := articles.Group("", auth, requireAdmin)
		{
			admin.POST("", apiLimiter, h.Create)
		}
	}
}
