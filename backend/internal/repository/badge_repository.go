package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/model"
)

// BadgeRepository 徽章仓储接口。
type BadgeRepository interface {
	Create(badge *model.Badge) error
	FindByUserAndType(userID uint64, badgeType string) (*model.Badge, error)
	ListByUserID(userID uint64) ([]model.Badge, error)
	CountByUserID(userID uint64) (int64, error)
}

type badgeRepository struct {
	db *gorm.DB
}

// NewBadgeRepository 构造徽章仓储。
func NewBadgeRepository(db *gorm.DB) BadgeRepository {
	return &badgeRepository{db: db}
}

func (r *badgeRepository) Create(badge *model.Badge) error {
	if err := r.db.Create(badge).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create badge type %s for user %d: %w", badge.Type, badge.UserID, ErrDuplicate)
		}
		return fmt.Errorf("create badge type %s for user %d: %w", badge.Type, badge.UserID, err)
	}
	return nil
}

func (r *badgeRepository) FindByUserAndType(userID uint64, badgeType string) (*model.Badge, error) {
	var badge model.Badge
	if err := r.db.Where("user_id = ? AND type = ?", userID, badgeType).First(&badge).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find badge %s of user %d: %w", badgeType, userID, ErrNotFound)
		}
		return nil, fmt.Errorf("find badge %s of user %d: %w", badgeType, userID, err)
	}
	return &badge, nil
}

func (r *badgeRepository) ListByUserID(userID uint64) ([]model.Badge, error) {
	var badges []model.Badge
	if err := r.db.Where("user_id = ?", userID).Order("earned_at DESC").Find(&badges).Error; err != nil {
		return nil, fmt.Errorf("list badges of user %d: %w", userID, err)
	}
	return badges, nil
}

func (r *badgeRepository) CountByUserID(userID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Badge{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count badges of user %d: %w", userID, err)
	}
	return total, nil
}
