package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wjecoffeetaste/wjecoffeetaste/internal/constants"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/dto"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/middleware"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/service"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/util"
)

// BeanFavoriteHandler exposes bean favorite endpoints.
type BeanFavoriteHandler struct {
	svc    *service.BeanFavoriteService
	logger *slog.Logger
}

// NewBeanFavoriteHandler creates a BeanFavoriteHandler.
func NewBeanFavoriteHandler(svc *service.BeanFavoriteService, logger *slog.Logger) *BeanFavoriteHandler {
	return &BeanFavoriteHandler{svc: svc, logger: logger}
}

// List handles GET /beans/favorites — the current user's own favorites only.
func (h *BeanFavoriteHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	origin := c.Query("origin")
	process := c.Query("process")
	keyword := c.Query("keyword")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}
	userID := middleware.GetUserID(c)
	items, total, err := h.svc.ListByUser(userID, origin, process, keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	ids := make([]uint, 0, len(items))
	for _, b := range items {
		ids = append(ids, b.ID)
	}
	counts, _ := h.svc.CountByBeans(ids)
	list := make([]dto.BeanItemResponse, 0, len(items))
	for _, b := range items {
		list = append(list, dto.BeanItemResponse{
			CoffeeBean:    b,
			FavoriteCount: counts[b.ID],
			IsFavorite:    true,
		})
	}
	c.JSON(http.StatusOK, dto.OK(dto.PageData{List: list, Total: total, Page: page, Size: pageSize}))
}

// Favorite handles POST /beans/:id/favorite.
func (h *BeanFavoriteHandler) Favorite(c *gin.Context) {
	beanID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid bean id"))
		return
	}
	f, err := h.svc.Favorite(middleware.GetUserID(c), uint(beanID))
	if err != nil {
		c.Error(err)
		return
	}
	count, _ := h.svc.CountByBean(uint(beanID))
	c.JSON(http.StatusCreated, dto.OK(gin.H{"favorite": f, "favorite_count": count}))
}

// Unfavorite handles DELETE /beans/:id/favorite.
func (h *BeanFavoriteHandler) Unfavorite(c *gin.Context) {
	beanID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(util.NewAppError(http.StatusBadRequest, constants.CodeBadRequest, "invalid bean id"))
		return
	}
	if err := h.svc.Unfavorite(middleware.GetUserID(c), uint(beanID)); err != nil {
		c.Error(err)
		return
	}
	count, _ := h.svc.CountByBean(uint(beanID))
	c.JSON(http.StatusOK, dto.OK(gin.H{"unfavorited": true, "favorite_count": count}))
}
