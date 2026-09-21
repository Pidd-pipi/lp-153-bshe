package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/model"
)

// TimeCapsuleRepository 时光胶囊仓储接口。
type TimeCapsuleRepository interface {
	Create(capsule *model.TimeCapsule) error
	FindByID(id uint64) (*model.TimeCapsule, error)
	Delete(id uint64) error
	ListByUserID(userID uint64, offset, limit int) ([]model.TimeCapsule, error)
	CountByUserID(userID uint64) (int64, error)
	UnlockDue(now time.Time) ([]model.TimeCapsule, error)
	Update(capsule *model.TimeCapsule) error
}

type timeCapsuleRepository struct {
	db *gorm.DB
}

// NewTimeCapsuleRepository 构造胶囊仓储。
func NewTimeCapsuleRepository(db *gorm.DB) TimeCapsuleRepository {
	return &timeCapsuleRepository{db: db}
}

func (r *timeCapsuleRepository) Create(capsule *model.TimeCapsule) error {
	if err := r.db.Create(capsule).Error; err != nil {
		return fmt.Errorf("create capsule: %w", err)
	}
	return nil
}

func (r *timeCapsuleRepository) FindByID(id uint64) (*model.TimeCapsule, error) {
	var capsule model.TimeCapsule
	if err := r.db.First(&capsule, id).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find capsule by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find capsule by id %d: %w", id, err)
	}
	return &capsule, nil
}

func (r *timeCapsuleRepository) Delete(id uint64) error {
	if err := r.db.Delete(&model.TimeCapsule{}, id).Error; err != nil {
		return fmt.Errorf("delete capsule %d: %w", id, err)
	}
	return nil
}

func (r *timeCapsuleRepository) ListByUserID(userID uint64, offset, limit int) ([]model.TimeCapsule, error) {
	var items []model.TimeCapsule
	if err := r.db.Where("user_id = ?", userID).Order("unlock_at ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list capsules by user %d: %w", userID, err)
	}
	return items, nil
}

func (r *timeCapsuleRepository) CountByUserID(userID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&model.TimeCapsule{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count capsules by user %d: %w", userID, err)
	}
	return total, nil
}

// UnlockDue 查询已到解锁时间但仍为 locked 的胶囊。
func (r *timeCapsuleRepository) UnlockDue(now time.Time) ([]model.TimeCapsule, error) {
	var items []model.TimeCapsule
	if err := r.db.Where("status = ? AND unlock_at <= ?", "locked", now).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list due capsules: %w", err)
	}
	return items, nil
}

func (r *timeCapsuleRepository) Update(capsule *model.TimeCapsule) error {
	if err := r.db.Save(capsule).Error; err != nil {
		return fmt.Errorf("update capsule %d: %w", capsule.ID, err)
	}
	return nil
}
