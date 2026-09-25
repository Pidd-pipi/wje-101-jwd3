package dto

import "github.com/wjecoffeetaste/wjecoffeetaste/internal/model"

// BeanCreateRequest creates/updates a coffee bean.
type BeanCreateRequest struct {
	Name          string `json:"name" binding:"required,max=128"`
	Origin        string `json:"origin" binding:"omitempty,max=128"`
	ProcessMethod string `json:"process_method" binding:"required"`
	FlavorTags    string `json:"flavor_tags"`
	Description   string `json:"description"`
}

// BeanItem decorates a coffee bean with favorite stats for the viewer.
type BeanItem struct {
	model.CoffeeBean
	FavoriteCount int64 `json:"favorite_count"`
	IsFavorited   bool  `json:"is_favorited"`
}
