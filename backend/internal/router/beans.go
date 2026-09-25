package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wjecoffeetaste/wjecoffeetaste/internal/config"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/handler"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/middleware"
)

func registerBeanRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.BeanHandler, limiter *middleware.RateLimiter) {
	beans := v1.Group("/beans")
	beans.GET("", middleware.AuthOptional(cfg), h.List)

	auth := beans.Group("", middleware.AuthRequired(cfg))
	auth.GET("/favorites", h.MyFavorites)
	auth.POST("/:id/favorite", limiter.Limit(), h.Favorite)
	auth.DELETE("/:id/favorite", h.Unfavorite)

	admin := beans.Group("", middleware.AuthRequired(cfg), middleware.RequireRole("admin"))
	admin.POST("", limiter.Limit(), h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}
