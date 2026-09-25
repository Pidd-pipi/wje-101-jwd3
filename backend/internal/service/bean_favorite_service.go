package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/wjecoffeetaste/wjecoffeetaste/internal/constants"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/dto"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/model"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/repository"
	"github.com/wjecoffeetaste/wjecoffeetaste/internal/util"
)

// BeanFavoriteService handles coffee bean favorites.
type BeanFavoriteService struct {
	repo     *repository.BeanFavoriteRepository
	beanRepo *repository.CoffeeBeanRepository
	logger   *slog.Logger
}

// NewBeanFavoriteService creates a BeanFavoriteService.
func NewBeanFavoriteService(repo *repository.BeanFavoriteRepository, beanRepo *repository.CoffeeBeanRepository, logger *slog.Logger) *BeanFavoriteService {
	return &BeanFavoriteService{repo: repo, beanRepo: beanRepo, logger: logger}
}

// Favorite adds a bean to the user's favorites (idempotent).
func (s *BeanFavoriteService) Favorite(userID, beanID uint) (*model.BeanFavorite, error) {
	if _, err := s.beanRepo.FindByID(beanID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("CoffeeBean[id=%d] not found", beanID))
		}
		return nil, fmt.Errorf("bean favorite find bean: %w", err)
	}
	if existing, err := s.repo.Find(userID, beanID); err == nil {
		return existing, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("bean favorite find: %w", err)
	}
	f := &model.BeanFavorite{UserID: userID, BeanID: beanID}
	if err := s.repo.Create(f); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			if existing, findErr := s.repo.Find(userID, beanID); findErr == nil {
				return existing, nil
			}
		}
		s.logger.Error(fmt.Sprintf(constants.LogBeanFavoriteFailed, beanID), "error", err)
		return nil, fmt.Errorf("bean favorite: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogBeanFavoriteSuccess, beanID), "user_id", userID)
	return f, nil
}

// Unfavorite removes a bean from the user's favorites (idempotent).
func (s *BeanFavoriteService) Unfavorite(userID, beanID uint) error {
	affected, err := s.repo.DeleteByUserBean(userID, beanID)
	if err != nil {
		return fmt.Errorf("bean unfavorite: %w", err)
	}
	if affected > 0 {
		s.logger.Info(fmt.Sprintf(constants.LogBeanUnfavoriteSuccess, beanID), "user_id", userID)
	}
	return nil
}

// ListByUser returns the user's favorited beans decorated with favorite stats.
func (s *BeanFavoriteService) ListByUser(userID uint, origin, process, keyword string, page, pageSize int) ([]dto.BeanItem, int64, error) {
	items, total, err := s.repo.ListBeansByUser(userID, origin, process, keyword, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("bean favorite list: %w", err)
	}
	decorated, err := s.decorate(items, userID)
	if err != nil {
		return nil, 0, err
	}
	return decorated, total, nil
}

// Decorate attaches favorite counts and the viewer's favorite flag to bean items.
// When viewerID is 0 (anonymous), every item reports is_favorited=false.
func (s *BeanFavoriteService) Decorate(items []model.CoffeeBean, viewerID uint) ([]dto.BeanItem, error) {
	return s.decorate(items, viewerID)
}

func (s *BeanFavoriteService) decorate(items []model.CoffeeBean, viewerID uint) ([]dto.BeanItem, error) {
	ids := make([]uint, 0, len(items))
	for _, b := range items {
		ids = append(ids, b.ID)
	}
	counts, err := s.repo.CountsByBeans(ids)
	if err != nil {
		return nil, fmt.Errorf("bean favorite counts: %w", err)
	}
	favorited, err := s.repo.FavoriteBeanIDs(viewerID, ids)
	if err != nil {
		return nil, fmt.Errorf("bean favorite ids: %w", err)
	}
	result := make([]dto.BeanItem, 0, len(items))
	for _, b := range items {
		result = append(result, dto.BeanItem{
			CoffeeBean:    b,
			FavoriteCount: counts[b.ID],
			IsFavorited:   favorited[b.ID],
		})
	}
	return result, nil
}
