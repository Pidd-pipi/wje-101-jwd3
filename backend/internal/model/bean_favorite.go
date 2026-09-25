package model

import "time"

// BeanFavorite links a user to a coffee bean (unique pair).
type BeanFavorite struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index:idx_bean_favorite_user_bean,unique;not null" json:"user_id"`
	BeanID    uint       `gorm:"index:idx_bean_favorite_user_bean,unique;not null" json:"bean_id"`
	User      *User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Bean      *CoffeeBean `gorm:"foreignKey:BeanID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
}
