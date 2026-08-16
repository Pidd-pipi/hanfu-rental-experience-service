// Package router assembles the Gin engine and all route groups for hanfu-rental.
package router

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lp/hanfu-rental/internal/config"
	"github.com/lp/hanfu-rental/internal/constants"
	"github.com/lp/hanfu-rental/internal/handler"
	"github.com/lp/hanfu-rental/internal/middleware"
	"github.com/lp/hanfu-rental/internal/repository"
	"github.com/lp/hanfu-rental/internal/service"
	"github.com/lp/hanfu-rental/internal/util"
	"gorm.io/gorm"
)

// New builds the Gin engine with all dependencies wired.
func New(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID(logger))
	r.Use(middleware.AccessLog(logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.CORSOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-Request-ID"},
	}))
	r.Use(middleware.ErrorHandler())

	r.GET("/healthz", func(c *gin.Context) { util.OK(c, gin.H{"status": "ok"}) })

	// repositories
	userRepo := repository.NewUserRepository(db)
	hanfuRepo := repository.NewHanfuRepository(db)
	catalogRepo := repository.NewCatalogRepository(db)
	orderRepo := repository.NewRentalOrderRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	cardRepo := repository.NewMemberCardRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// services
	userSvc := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpireHours, logger)
	hanfuSvc := service.NewHanfuService(hanfuRepo, logger)
	orderSvc := service.NewRentalOrderService(orderRepo, hanfuRepo, catalogRepo, cardRepo, userRepo, logger)
	activitySvc := service.NewActivityService(activityRepo, logger)
	cardSvc := service.NewMemberCardService(cardRepo, logger)
	articleSvc := service.NewArticleService(articleRepo, logger)

	// handlers
	userH := handler.NewUserHandler(userSvc, logger)
	hanfuH := handler.NewHanfuHandler(hanfuSvc, logger)
	catalogH := handler.NewCatalogHandler(catalogRepo, logger)
	orderH := handler.NewRentalOrderHandler(orderSvc, userSvc, logger)
	activityH := handler.NewActivityHandler(activitySvc, userSvc, logger)
	cardH := handler.NewMemberCardHandler(cardSvc, logger)
	articleH := handler.NewArticleHandler(articleSvc, logger)

	auth := middleware.AuthRequired(cfg.JWTSecret, logger)
	requireAdmin := middleware.RequireRole(logger, constants.UserRoleAdmin)
	loginLimiter := middleware.RateLimit(middleware.NewRateLimiter(cfg.LoginRateLimit, time.Minute), logger)
	apiLimiter := middleware.RateLimit(middleware.NewRateLimiter(cfg.RateLimitPerMin, time.Minute), logger)

	v1 := r.Group("/api/v1")
	{
		RegisterUserRoutes(v1, userH, auth, requireAdmin, loginLimiter, apiLimiter)
		RegisterHanfuRoutes(v1, hanfuH, auth, requireAdmin, apiLimiter)
		RegisterRentalOrderRoutes(v1, orderH, auth, requireAdmin, apiLimiter)
		RegisterActivityRoutes(v1, activityH, auth, requireAdmin, apiLimiter)
		RegisterMemberCardRoutes(v1, cardH, auth, apiLimiter)
		RegisterArticleRoutes(v1, articleH, auth, requireAdmin, apiLimiter)
		// makeup packages & photographers
		v1.GET("/makeup-packages", apiLimiter, catalogH.ListMakeupPackages)
		v1.GET("/photographers", apiLimiter, catalogH.ListPhotographers)
	}
	return r
}
