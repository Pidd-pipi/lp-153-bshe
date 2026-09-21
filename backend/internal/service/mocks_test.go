package service

import (
	"log/slog"
	"os"

	"gorm.io/gorm"

	"github.com/wishwall/wishwall/internal/config"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testConfig() *config.Config {
	return &config.Config{JWTSecret: "test-secret-for-unit", JWTExpireHours: 1, RateLimitPerMin: 1000}
}

// ---- mock UserRepository ----
type mockUserRepo struct {
	createFn         func(user *model.User) error
	findByIDFn       func(id uint64) (*model.User, error)
	findByUsernameFn func(username string) (*model.User, error)
	findByEmailFn    func(email string) (*model.User, error)
	findByAccountFn  func(account string) (*model.User, error)
	updateFn         func(user *model.User) error
}

func (m *mockUserRepo) Create(user *model.User) error           { return m.createFn(user) }
func (m *mockUserRepo) FindByID(id uint64) (*model.User, error) { return m.findByIDFn(id) }
func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	return m.findByUsernameFn(username)
}
func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) { return m.findByEmailFn(email) }
func (m *mockUserRepo) FindByAccount(account string) (*model.User, error) {
	return m.findByAccountFn(account)
}
func (m *mockUserRepo) Update(user *model.User) error { return m.updateFn(user) }

// ---- mock AuditService ----
type mockAudit struct {
	recordFn func(log *model.AuditLog) error
	listFn   func(q dto.AuditLogQuery) (*dto.PageResult, error)
}

func (m *mockAudit) Record(log *model.AuditLog) error {
	if m.recordFn == nil {
		return nil
	}
	return m.recordFn(log)
}
func (m *mockAudit) List(q dto.AuditLogQuery) (*dto.PageResult, error) { return m.listFn(q) }

// ---- mock BadgeService ----
type mockBadge struct {
	grantFirstWishFn     func(userID uint64) error
	grantFirstClaimFn    func(userID uint64) error
	grantFirstBlessingFn func(userID uint64) error
	grantCompletionFn    func(userID uint64) error
	listMineFn           func(userID uint64) ([]dto.BadgeResponse, error)
	leaderboardFn        func(limit int) ([]model.FulfillerStat, error)
}

func (m *mockBadge) GrantFirstWish(userID uint64) error                  { return m.grantFirstWishFn(userID) }
func (m *mockBadge) GrantFirstClaim(userID uint64) error                 { return m.grantFirstClaimFn(userID) }
func (m *mockBadge) GrantFirstBlessing(userID uint64) error              { return m.grantFirstBlessingFn(userID) }
func (m *mockBadge) GrantCompletionBadges(userID uint64) error           { return m.grantCompletionFn(userID) }
func (m *mockBadge) ListMine(userID uint64) ([]dto.BadgeResponse, error) { return m.listMineFn(userID) }
func (m *mockBadge) Leaderboard(limit int) ([]model.FulfillerStat, error) {
	return m.leaderboardFn(limit)
}

// ---- mock WishRepository ----
type mockWishRepo struct {
	createFn            func(wish *model.Wish) error
	findByIDFn          func(id uint64) (*model.Wish, error)
	findByIDForUpdateFn func(tx *gorm.DB, id uint64) (*model.Wish, error)
	updateFn            func(wish *model.Wish) error
	updateWithTxFn      func(tx *gorm.DB, wish *model.Wish) error
	deleteFn            func(id uint64) error
	listFn              func(filters map[string]any, sort string, offset, limit int) ([]model.Wish, error)
	countFn             func(filters map[string]any) (int64, error)
	listByUserFn        func(userID uint64, offset, limit int) ([]model.Wish, error)
	countByUserFn       func(userID uint64) (int64, error)
	incrementLikesFn    func(wishID uint64) error
	countCompletedFn    func(userID uint64) (int64, error)
	listCompletedFn     func(offset, limit int) ([]model.Wish, error)
	topFulfillersFn     func(limit int) ([]model.FulfillerStat, error)
}

func (m *mockWishRepo) Create(wish *model.Wish) error           { return m.createFn(wish) }
func (m *mockWishRepo) FindByID(id uint64) (*model.Wish, error) { return m.findByIDFn(id) }
func (m *mockWishRepo) FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.Wish, error) {
	return m.findByIDForUpdateFn(tx, id)
}
func (m *mockWishRepo) Update(wish *model.Wish) error { return m.updateFn(wish) }
func (m *mockWishRepo) UpdateWithTx(tx *gorm.DB, wish *model.Wish) error {
	return m.updateWithTxFn(tx, wish)
}
func (m *mockWishRepo) Delete(id uint64) error { return m.deleteFn(id) }
func (m *mockWishRepo) List(filters map[string]any, sort string, offset, limit int) ([]model.Wish, error) {
	return m.listFn(filters, sort, offset, limit)
}
func (m *mockWishRepo) Count(filters map[string]any) (int64, error) { return m.countFn(filters) }
func (m *mockWishRepo) ListByUserID(userID uint64, offset, limit int) ([]model.Wish, error) {
	return m.listByUserFn(userID, offset, limit)
}
func (m *mockWishRepo) CountByUserID(userID uint64) (int64, error) { return m.countByUserFn(userID) }
func (m *mockWishRepo) IncrementLikes(wishID uint64) error         { return m.incrementLikesFn(wishID) }
func (m *mockWishRepo) CountCompletedByFulfiller(userID uint64) (int64, error) {
	return m.countCompletedFn(userID)
}
func (m *mockWishRepo) ListCompletedStories(offset, limit int) ([]model.Wish, error) {
	return m.listCompletedFn(offset, limit)
}
func (m *mockWishRepo) TopFulfillers(limit int) ([]model.FulfillerStat, error) {
	return m.topFulfillersFn(limit)
}

// ---- mock WishClaimRepository ----
type mockClaimRepo struct {
	createFn            func(claim *model.WishClaim) error
	createWithTxFn      func(tx *gorm.DB, claim *model.WishClaim) error
	findByIDFn          func(id uint64) (*model.WishClaim, error)
	findByIDForUpdateFn func(tx *gorm.DB, id uint64) (*model.WishClaim, error)
	findByWishFn        func(wishID uint64) (*model.WishClaim, error)
	updateFn            func(claim *model.WishClaim) error
	updateWithTxFn      func(tx *gorm.DB, claim *model.WishClaim) error
	listByUserFn        func(userID uint64, offset, limit int) ([]model.WishClaim, error)
	countByUserFn       func(userID uint64) (int64, error)
	listByWishFn        func(wishIDs []uint64) ([]model.WishClaim, error)
}

func (m *mockClaimRepo) Create(claim *model.WishClaim) error { return m.createFn(claim) }
func (m *mockClaimRepo) CreateWithTx(tx *gorm.DB, claim *model.WishClaim) error {
	return m.createWithTxFn(tx, claim)
}
func (m *mockClaimRepo) FindByID(id uint64) (*model.WishClaim, error) { return m.findByIDFn(id) }
func (m *mockClaimRepo) FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.WishClaim, error) {
	return m.findByIDForUpdateFn(tx, id)
}
func (m *mockClaimRepo) FindByWishID(wishID uint64) (*model.WishClaim, error) {
	return m.findByWishFn(wishID)
}
func (m *mockClaimRepo) Update(claim *model.WishClaim) error { return m.updateFn(claim) }
func (m *mockClaimRepo) UpdateWithTx(tx *gorm.DB, claim *model.WishClaim) error {
	return m.updateWithTxFn(tx, claim)
}
func (m *mockClaimRepo) ListByUserID(userID uint64, offset, limit int) ([]model.WishClaim, error) {
	return m.listByUserFn(userID, offset, limit)
}
func (m *mockClaimRepo) CountByUserID(userID uint64) (int64, error) { return m.countByUserFn(userID) }
func (m *mockClaimRepo) ListByWishIDs(wishIDs []uint64) ([]model.WishClaim, error) {
	return m.listByWishFn(wishIDs)
}

// ---- mock BlessingRepository ----
type mockBlessRepo struct {
	createFn      func(blessing *model.Blessing) error
	listByWishFn  func(wishID uint64, offset, limit int) ([]model.Blessing, error)
	countByWishFn func(wishID uint64) (int64, error)
	countByUserFn func(userID uint64) (int64, error)
}

func (m *mockBlessRepo) Create(b *model.Blessing) error { return m.createFn(b) }
func (m *mockBlessRepo) ListByWishID(wishID uint64, offset, limit int) ([]model.Blessing, error) {
	return m.listByWishFn(wishID, offset, limit)
}
func (m *mockBlessRepo) CountByWishID(wishID uint64) (int64, error) { return m.countByWishFn(wishID) }
func (m *mockBlessRepo) CountByUserID(userID uint64) (int64, error) { return m.countByUserFn(userID) }

// ---- mock TxManager ----
type mockTx struct {
	fn func(fn func(tx *gorm.DB) error) error
}

func (m *mockTx) Transaction(fn func(tx *gorm.DB) error) error {
	if m.fn == nil {
		return fn(nil)
	}
	return m.fn(fn)
}

// notFoundErr 快捷构造仓储未找到错误。
func notFoundErr() error { return repository.ErrNotFound }
