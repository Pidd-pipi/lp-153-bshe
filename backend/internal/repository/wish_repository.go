package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/model"
)

// WishRepository 心愿仓储接口。
type WishRepository interface {
	Create(wish *model.Wish) error
	FindByID(id uint64) (*model.Wish, error)
	FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.Wish, error)
	Update(wish *model.Wish) error
	UpdateWithTx(tx *gorm.DB, wish *model.Wish) error
	Delete(id uint64) error
	List(filters map[string]any, sort string, offset, limit int) ([]model.Wish, error)
	Count(filters map[string]any) (int64, error)
	ListByUserID(userID uint64, offset, limit int) ([]model.Wish, error)
	CountByUserID(userID uint64) (int64, error)
	IncrementLikes(wishID uint64) error
	CountCompletedByFulfiller(userID uint64) (int64, error)
	ListCompletedStories(offset, limit int) ([]model.Wish, error)
	TopFulfillers(limit int) ([]model.FulfillerStat, error)
}

type wishRepository struct {
	db *gorm.DB
}

// NewWishRepository 构造心愿仓储。
func NewWishRepository(db *gorm.DB) WishRepository {
	return &wishRepository{db: db}
}

func (r *wishRepository) Create(wish *model.Wish) error {
	if err := r.db.Create(wish).Error; err != nil {
		return fmt.Errorf("create wish: %w", err)
	}
	return nil
}

func (r *wishRepository) FindByID(id uint64) (*model.Wish, error) {
	var wish model.Wish
	if err := r.db.First(&wish, id).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find wish by id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find wish by id %d: %w", id, err)
	}
	return &wish, nil
}

// FindByIDForUpdate 行级锁读取，用于并发认领场景。
func (r *wishRepository) FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.Wish, error) {
	var wish model.Wish
	if err := tx.Clauses(gormclauseLock()).First(&wish, id).Error; err != nil {
		if isRecordNotFound(err) {
			return nil, fmt.Errorf("find wish by id %d for update: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("find wish by id %d for update: %w", id, err)
	}
	return &wish, nil
}

func (r *wishRepository) Update(wish *model.Wish) error {
	if err := r.db.Save(wish).Error; err != nil {
		return fmt.Errorf("update wish %d: %w", wish.ID, err)
	}
	return nil
}

func (r *wishRepository) UpdateWithTx(tx *gorm.DB, wish *model.Wish) error {
	if err := tx.Save(wish).Error; err != nil {
		return fmt.Errorf("update wish %d with tx: %w", wish.ID, err)
	}
	return nil
}

func (r *wishRepository) Delete(id uint64) error {
	if err := r.db.Delete(&model.Wish{}, id).Error; err != nil {
		return fmt.Errorf("delete wish %d: %w", id, err)
	}
	return nil
}

func (r *wishRepository) List(filters map[string]any, sort string, offset, limit int) ([]model.Wish, error) {
	var wishes []model.Wish
	q := r.db.Model(&model.Wish{})
	q = applyWishFilters(q, filters)
	if sort == "hot" {
		q = q.Order("likes_count DESC, created_at DESC")
	} else {
		q = q.Order("created_at DESC")
	}
	if err := q.Offset(offset).Limit(limit).Find(&wishes).Error; err != nil {
		return nil, fmt.Errorf("list wishes: %w", err)
	}
	return wishes, nil
}

func (r *wishRepository) Count(filters map[string]any) (int64, error) {
	var total int64
	q := r.db.Model(&model.Wish{})
	q = applyWishFilters(q, filters)
	if err := q.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count wishes: %w", err)
	}
	return total, nil
}

func (r *wishRepository) ListByUserID(userID uint64, offset, limit int) ([]model.Wish, error) {
	var wishes []model.Wish
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&wishes).Error; err != nil {
		return nil, fmt.Errorf("list wishes by user %d: %w", userID, err)
	}
	return wishes, nil
}

func (r *wishRepository) CountByUserID(userID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Wish{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count wishes by user %d: %w", userID, err)
	}
	return total, nil
}

func (r *wishRepository) IncrementLikes(wishID uint64) error {
	if err := r.db.Model(&model.Wish{}).Where("id = ?", wishID).
		UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error; err != nil {
		return fmt.Errorf("increment likes of wish %d: %w", wishID, err)
	}
	return nil
}

// CountCompletedByFulfiller 统计某用户作为圆梦人完成的心愿数（成就/排行榜复用）。
func (r *wishRepository) CountCompletedByFulfiller(userID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Wish{}).
		Joins("JOIN wish_claims ON wish_claims.wish_id = wishes.id").
		Where("wish_claims.user_id = ? AND wishes.status = ?", userID, "completed").
		Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count completed by fulfiller %d: %w", userID, err)
	}
	return total, nil
}

func (r *wishRepository) TopFulfillers(limit int) ([]model.FulfillerStat, error) {
	var rows []model.FulfillerStat
	if err := r.db.Model(&model.WishClaim{}).
		Select("wish_claims.user_id AS user_id, users.nickname AS nickname, users.avatar AS avatar, COUNT(wishes.id) AS completed_count").
		Joins("JOIN wishes ON wishes.id = wish_claims.wish_id AND wishes.status = ?", "completed").
		Joins("JOIN users ON users.id = wish_claims.user_id").
		Group("wish_claims.user_id, users.nickname, users.avatar").
		Order("completed_count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("top fulfillers: %w", err)
	}
	return rows, nil
}

func (r *wishRepository) ListCompletedStories(offset, limit int) ([]model.Wish, error) {
	var wishes []model.Wish
	if err := r.db.Where("status = ?", "completed").
		Order("updated_at DESC").Offset(offset).Limit(limit).Find(&wishes).Error; err != nil {
		return nil, fmt.Errorf("list completed stories: %w", err)
	}
	return wishes, nil
}
