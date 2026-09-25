package repository

import (
	"gorm.io/gorm"

	"github.com/wjecoffeetaste/wjecoffeetaste/internal/model"
)

// BeanFavoriteRepository handles bean favorite persistence.
type BeanFavoriteRepository struct{ db *gorm.DB }

// NewBeanFavoriteRepository creates the repository.
func NewBeanFavoriteRepository(db *gorm.DB) *BeanFavoriteRepository {
	return &BeanFavoriteRepository{db: db}
}

// Create inserts a favorite.
func (r *BeanFavoriteRepository) Create(f *model.BeanFavorite) error {
	return translate(r.db.Create(f).Error)
}

// Find locates a favorite by user and bean.
func (r *BeanFavoriteRepository) Find(userID, beanID uint) (*model.BeanFavorite, error) {
	var f model.BeanFavorite
	if err := translate(r.db.Where("user_id = ? AND bean_id = ?", userID, beanID).First(&f).Error); err != nil {
		return nil, err
	}
	return &f, nil
}

// DeleteByUserBean removes a favorite by user and bean.
func (r *BeanFavoriteRepository) DeleteByUserBean(userID, beanID uint) error {
	return r.db.Where("user_id = ? AND bean_id = ?", userID, beanID).Delete(&model.BeanFavorite{}).Error
}

// DeleteByBean removes all favorites of a bean (used when admin deletes it).
func (r *BeanFavoriteRepository) DeleteByBean(beanID uint) error {
	return r.db.Where("bean_id = ?", beanID).Delete(&model.BeanFavorite{}).Error
}

// FindExistingBeanIDs returns the subset of bean ids the user has favorited.
func (r *BeanFavoriteRepository) FindExistingBeanIDs(userID uint, beanIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if len(beanIDs) == 0 {
		return result, nil
	}
	var ids []uint
	if err := r.db.Model(&model.BeanFavorite{}).
		Where("user_id = ? AND bean_id IN ?", userID, beanIDs).
		Pluck("bean_id", &ids).Error; err != nil {
		return nil, err
	}
	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}

// CountByBeans returns favorite counts keyed by bean id.
func (r *BeanFavoriteRepository) CountByBeans(beanIDs []uint) (map[uint]int64, error) {
	result := make(map[uint]int64)
	if len(beanIDs) == 0 {
		return result, nil
	}
	type countRow struct {
		BeanID uint
		Total  int64
	}
	var rows []countRow
	if err := r.db.Model(&model.BeanFavorite{}).
		Select("bean_id AS bean_id, COUNT(*) AS total").
		Where("bean_id IN ?", beanIDs).
		Group("bean_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.BeanID] = row.Total
	}
	return result, nil
}

// CountByBean counts favorites of a bean.
func (r *BeanFavoriteRepository) CountByBean(beanID uint) (int64, error) {
	var total int64
	if err := r.db.Model(&model.BeanFavorite{}).Where("bean_id = ?", beanID).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// ListBeansByUser lists favorite beans of a user with the same filters as the
// bean library, returning beans and the total count.
func (r *BeanFavoriteRepository) ListBeansByUser(userID uint, origin, process, keyword string, page, pageSize int) ([]model.CoffeeBean, int64, error) {
	var items []model.CoffeeBean
	var total int64
	q := r.db.Model(&model.CoffeeBean{}).
		Joins("JOIN bean_favorites ON bean_favorites.bean_id = coffee_beans.id").
		Where("bean_favorites.user_id = ?", userID)
	if origin != "" {
		q = q.Where("coffee_beans.origin = ?", origin)
	}
	if process != "" {
		q = q.Where("coffee_beans.process_method = ?", process)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("coffee_beans.name LIKE ? OR coffee_beans.flavor_tags LIKE ?", like, like)
	}
	// the (user_id, bean_id) unique index guarantees no duplicated beans.
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("bean_favorites.created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
