package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wjecoffeetaste/wjecoffeetaste/internal/config"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/handler"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/middleware"
)

func registerBeanRoutes(v1 *gin.RouterGroup, cfg *config.Config, h *handler.BeanHandler, fh *handler.BeanFavoriteHandler, limiter *middleware.RateLimiter) {
	beans := v1.Group("/beans")
	// Public list: honors an optional token so cards can show favorite state.
	beans.GET("", middleware.AuthOptional(cfg), h.List)
	beans.GET("/favorites", middleware.AuthRequired(cfg), fh.List)
	auth := beans.Group("", middleware.AuthRequired(cfg))
	auth.POST("/:id/favorite", limiter.Limit(), fh.Favorite)
	auth.DELETE("/:id/favorite", fh.Unfavorite)
	admin := beans.Group("", middleware.AuthRequired(cfg), middleware.RequireRole("admin"))
	admin.POST("", limiter.Limit(), h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}
