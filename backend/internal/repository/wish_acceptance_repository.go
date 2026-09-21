package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/model"
)

// WishAcceptanceRepository 心愿验收仓储接口。
type WishAcceptanceRepository interface {
	CreateWithTx(tx *gorm.DB, acc *model.WishAcceptance) error
	FindByID(id uint64) (*model.WishAcceptance, error)
	FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.WishAcceptance, error)
	FindByWishID(wishID uint64) (*model.WishAcceptance, error)
	FindByWishIDForUpdate(tx *gorm.DB, wishID uint64) (*model.WishAcceptance, error)
	UpdateWithTx(tx *gorm.DB, acc *model.WishAcceptance) error
}

type wishAcceptanceRepository struct {
	db *gorm.DB
}

// NewWishAcceptanceRepository 构造验收仓储。
func NewWishAcceptanceRepository(db *gorm.DB) WishAcceptanceRepository {
	return &wishAcceptanceRepository{db: db}
}

func (r *wishAcceptanceRepository) CreateWithTx(tx *gorm.DB, acc *model.WishAcceptance) error {
	if err := tx.Create(acc).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("create acceptance on wish %d with tx: %w", acc.WishID, ErrConflict)
		}
		return fmt.Errorf("create acceptance on wish %d with tx: %w", acc.WishID, err)
	}
	return nil
}

func (r *wishAcceptanceRepository) FindByID(id uint64) (*model.WishAcceptance, error) {
	var acc model.WishAcceptance
	if err := r.db.First(&acc, id).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find acceptance by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find acceptance by id %d: %w", id, err)
	}
	return &acc, nil
}

// FindByIDForUpdate 行级锁读取，用于并发验收（确认/驳回只能成功一次）。
func (r *wishAcceptanceRepository) FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.WishAcceptance, error) {
	var acc model.WishAcceptance
	if err := tx.Clauses(gormclauseLock()).First(&acc, id).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find acceptance by id %d for update: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find acceptance by id %d for update: %w", id, err)
	}
	return &acc, nil
}

func (r *wishAcceptanceRepository) FindByWishID(wishID uint64) (*model.WishAcceptance, error) {
	var acc model.WishAcceptance
	if err := r.db.Where("wish_id = ?", wishID).First(&acc).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find acceptance by wish %d: %w", wishID, ErrNotFound)
		}
		return nil, fmt.Errorf("find acceptance by wish %d: %w", wishID, err)
	}
	return &acc, nil
}

// FindByWishIDForUpdate 行级锁读取，用于送交/重新送交与验收互斥。
func (r *wishAcceptanceRepository) FindByWishIDForUpdate(tx *gorm.DB, wishID uint64) (*model.WishAcceptance, error) {
	var acc model.WishAcceptance
	if err := tx.Clauses(gormclauseLock()).Where("wish_id = ?", wishID).First(&acc).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find acceptance by wish %d for update: %w", wishID, ErrNotFound)
		}
		return nil, fmt.Errorf("find acceptance by wish %d for update: %w", wishID, err)
	}
	return &acc, nil
}

func (r *wishAcceptanceRepository) UpdateWithTx(tx *gorm.DB, acc *model.WishAcceptance) error {
	if err := tx.Save(acc).Error; err != nil {
		return fmt.Errorf("update acceptance %d with tx: %w", acc.ID, err)
	}
	return nil
}
