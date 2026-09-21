package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/model"
)

// BlessingRepository 祝福留言仓储接口。
type BlessingRepository interface {
	Create(blessing *model.Blessing) error
	ListByWishID(wishID uint64, offset, limit int) ([]model.Blessing, error)
	CountByWishID(wishID uint64) (int64, error)
	CountByUserID(userID uint64) (int64, error)
}

type blessingRepository struct {
	db *gorm.DB
}

// NewBlessingRepository 构造祝福仓储。
func NewBlessingRepository(db *gorm.DB) BlessingRepository {
	return &blessingRepository{db: db}
}

func (r *blessingRepository) Create(blessing *model.Blessing) error {
	if err := r.db.Create(blessing).Error; err != nil {
		return fmt.Errorf("create blessing on wish %d: %w", blessing.WishID, err)
	}
	return nil
}

func (r *blessingRepository) ListByWishID(wishID uint64, offset, limit int) ([]model.Blessing, error) {
	var items []model.Blessing
	if err := r.db.Where("wish_id = ?", wishID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list blessings of wish %d: %w", wishID, err)
	}
	return items, nil
}

func (r *blessingRepository) CountByWishID(wishID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Blessing{}).Where("wish_id = ?", wishID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count blessings of wish %d: %w", wishID, err)
	}
	return total, nil
}

func (r *blessingRepository) CountByUserID(userID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Blessing{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count blessings by user %d: %w", userID, err)
	}
	return total, nil
}
