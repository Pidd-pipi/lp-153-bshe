package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/model"
)

// WishClaimRepository 心愿认领仓储接口。
type WishClaimRepository interface {
	Create(claim *model.WishClaim) error
	CreateWithTx(tx *gorm.DB, claim *model.WishClaim) error
	FindByID(id uint64) (*model.WishClaim, error)
	FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.WishClaim, error)
	FindByWishID(wishID uint64) (*model.WishClaim, error)
	Update(claim *model.WishClaim) error
	UpdateWithTx(tx *gorm.DB, claim *model.WishClaim) error
	ListByUserID(userID uint64, offset, limit int) ([]model.WishClaim, error)
	CountByUserID(userID uint64) (int64, error)
	ListByWishIDs(wishIDs []uint64) ([]model.WishClaim, error)
}

type wishClaimRepository struct {
	db *gorm.DB
}

// NewWishClaimRepository 构造认领仓储。
func NewWishClaimRepository(db *gorm.DB) WishClaimRepository {
	return &wishClaimRepository{db: db}
}

func (r *wishClaimRepository) Create(claim *model.WishClaim) error {
	if err := r.db.Create(claim).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create claim on wish %d: %w", claim.WishID, ErrConflict)
		}
		return fmt.Errorf("create claim on wish %d: %w", claim.WishID, err)
	}
	return nil
}

func (r *wishClaimRepository) CreateWithTx(tx *gorm.DB, claim *model.WishClaim) error {
	if err := tx.Create(claim).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create claim on wish %d with tx: %w", claim.WishID, ErrConflict)
		}
		return fmt.Errorf("create claim on wish %d with tx: %w", claim.WishID, err)
	}
	return nil
}

func (r *wishClaimRepository) FindByID(id uint64) (*model.WishClaim, error) {
	var claim model.WishClaim
	if err := r.db.First(&claim, id).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find claim by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find claim by id %d: %w", id, err)
	}
	return &claim, nil
}

// FindByIDForUpdate 事务内行级锁读取认领记录（发布者并发验收，只允许一次成功）。
func (r *wishClaimRepository) FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.WishClaim, error) {
	var claim model.WishClaim
	if err := tx.Clauses(gormclauseLock()).First(&claim, id).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find claim by id %d for update: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find claim by id %d for update: %w", id, err)
	}
	return &claim, nil
}

func (r *wishClaimRepository) FindByWishID(wishID uint64) (*model.WishClaim, error) {
	var claim model.WishClaim
	if err := r.db.Where("wish_id = ?", wishID).First(&claim).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find claim by wish %d: %w", wishID, ErrNotFound)
		}
		return nil, fmt.Errorf("find claim by wish %d: %w", wishID, err)
	}
	return &claim, nil
}

func (r *wishClaimRepository) Update(claim *model.WishClaim) error {
	if err := r.db.Save(claim).Error; err != nil {
		return fmt.Errorf("update claim %d: %w", claim.ID, err)
	}
	return nil
}

func (r *wishClaimRepository) UpdateWithTx(tx *gorm.DB, claim *model.WishClaim) error {
	if err := tx.Save(claim).Error; err != nil {
		return fmt.Errorf("update claim %d with tx: %w", claim.ID, err)
	}
	return nil
}

func (r *wishClaimRepository) ListByUserID(userID uint64, offset, limit int) ([]model.WishClaim, error) {
	var claims []model.WishClaim
	if err := r.db.Where("user_id = ?", userID).Order("updated_at DESC").Offset(offset).Limit(limit).Find(&claims).Error; err != nil {
		return nil, fmt.Errorf("list claims by user %d: %w", userID, err)
	}
	return claims, nil
}

func (r *wishClaimRepository) CountByUserID(userID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&model.WishClaim{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count claims by user %d: %w", userID, err)
	}
	return total, nil
}

// ListByWishIDs 批量取心愿的认领记录（列表页组装圆梦人信息复用）。
func (r *wishClaimRepository) ListByWishIDs(wishIDs []uint64) ([]model.WishClaim, error) {
	var claims []model.WishClaim
	if len(wishIDs) == 0 {
		return claims, nil
	}
	if err := r.db.Where("wish_id IN ?", wishIDs).Find(&claims).Error; err != nil {
		return nil, fmt.Errorf("list claims by wish ids: %w", err)
	}
	return claims, nil
}
