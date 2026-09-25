package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/wjecoffeetaste/wjecoffeetaste/internal/constants"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/model"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/repository"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/util"
)

// BeanFavoriteService handles user-to-bean favorite relations.
type BeanFavoriteService struct {
	repo     *repository.BeanFavoriteRepository
	beanRepo *repository.CoffeeBeanRepository
	logger   *slog.Logger
}

// NewBeanFavoriteService creates a BeanFavoriteService.
func NewBeanFavoriteService(repo *repository.BeanFavoriteRepository, beanRepo *repository.CoffeeBeanRepository, logger *slog.Logger) *BeanFavoriteService {
	return &BeanFavoriteService{repo: repo, beanRepo: beanRepo, logger: logger}
}

// Favorite marks a bean as favorite for the user.
func (s *BeanFavoriteService) Favorite(userID, beanID uint) (*model.BeanFavorite, error) {
	if _, err := s.beanRepo.FindByID(beanID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound,
				fmt.Sprintf("CoffeeBean[id=%d] not found", beanID))
		}
		return nil, fmt.Errorf("bean favorite find bean: %w", err)
	}
	f := &model.BeanFavorite{UserID: userID, BeanID: beanID}
	if err := s.repo.Create(f); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, util.NewAppError(409, constants.CodeConflict,
				fmt.Sprintf("BeanFavorite[user_id=%d bean_id=%d] failed: already favorited", userID, beanID))
		}
		s.logger.Error(fmt.Sprintf(constants.LogBeanFavoriteFailed, beanID), "user_id", userID, "error", err)
		return nil, fmt.Errorf("bean favorite: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogBeanFavoriteSuccess, beanID), "user_id", userID)
	return f, nil
}

// Unfavorite removes a favorite; it is idempotent so double cancel is a no-op.
func (s *BeanFavoriteService) Unfavorite(userID, beanID uint) error {
	if err := s.repo.DeleteByUserBean(userID, beanID); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogBeanUnfavoriteFailed, beanID), "user_id", userID, "error", err)
		return fmt.Errorf("bean unfavorite: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogBeanUnfavoriteSuccess, beanID), "user_id", userID)
	return nil
}

// ListByUser returns the user's favorite beans.
func (s *BeanFavoriteService) ListByUser(userID uint, origin, process, keyword string, page, pageSize int) ([]model.CoffeeBean, int64, error) {
	items, total, err := s.repo.ListBeansByUser(userID, origin, process, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("bean favorite list: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogBeanFavoriteListSuccess, userID), "total", total)
	return items, total, nil
}

// FindExistingBeanIDs returns which beans the user has favorited.
func (s *BeanFavoriteService) FindExistingBeanIDs(userID uint, beanIDs []uint) (map[uint]bool, error) {
	return s.repo.FindExistingBeanIDs(userID, beanIDs)
}

// CountByBeans returns favorite counts keyed by bean id.
func (s *BeanFavoriteService) CountByBeans(beanIDs []uint) (map[uint]int64, error) {
	return s.repo.CountByBeans(beanIDs)
}

// CountByBean returns the favorite count of one bean.
func (s *BeanFavoriteService) CountByBean(beanID uint) (int64, error) {
	return s.repo.CountByBean(beanID)
}
