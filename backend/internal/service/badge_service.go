package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

// BadgeService 成就徽章服务：按事件发放徽章、查询个人徽章与排行榜。
type BadgeService interface {
	GrantFirstWish(userID uint64) error
	GrantFirstClaim(userID uint64) error
	GrantFirstBlessing(userID uint64) error
	GrantCompletionBadges(userID uint64) error
	ListMine(userID uint64) ([]dto.BadgeResponse, error)
	Leaderboard(limit int) ([]model.FulfillerStat, error)
}

type badgeService struct {
	badge  repository.BadgeRepository
	wish   repository.WishRepository
	redis  *redis.Client
	logger *slog.Logger
}

// NewBadgeService 构造徽章服务。
func NewBadgeService(badge repository.BadgeRepository, wish repository.WishRepository, rdb *redis.Client, logger *slog.Logger) BadgeService {
	return &badgeService{badge: badge, wish: wish, redis: rdb, logger: logger}
}

// grant 幂等发放徽章：已存在则跳过。
func (s *badgeService) grant(userID uint64, badgeType, title, desc, icon string) error {
	if _, err := s.badge.FindByUserAndType(userID, badgeType); err == nil {
		return nil
	}
	badge := &model.Badge{
		UserID: userID, Type: badgeType, Title: title,
		Description: desc, Icon: icon, EarnedAt: time.Now(),
	}
	if err := s.badge.Create(badge); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil
		}
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogBadgeGranted, "user_id", userID, "badge_type", badgeType)
	return nil
}

func (s *badgeService) GrantFirstWish(userID uint64) error {
	total, err := s.wish.CountByUserID(userID)
	if err != nil {
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if total == 1 {
		return s.grant(userID, constants.BadgeTypeFirstWish, "首次许愿", "发布你的第一个心愿", "🌟")
	}
	return nil
}

func (s *badgeService) GrantFirstClaim(userID uint64) error {
	// 认领成功即触发（幂等判断在 grant 内）。
	return s.grant(userID, constants.BadgeTypeFirstClaim, "首次认领", "认领他人心愿成为圆梦人", "🤝")
}

func (s *badgeService) GrantFirstBlessing(userID uint64) error {
	return s.grant(userID, constants.BadgeTypeFirstBlessing, "首次祝福", "送出第一份祝福与鼓励", "💌")
}

func (s *badgeService) GrantCompletionBadges(userID uint64) error {
	total, err := s.wish.CountCompletedByFulfiller(userID)
	if err != nil {
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if total >= 10 {
		if err := s.grant(userID, constants.BadgeTypeTenCompletions, "十次圆梦", "作为圆梦人完成 10 个心愿", "🏆"); err != nil {
			return err
		}
	}
	if total >= 20 {
		if err := s.grant(userID, constants.BadgeTypeWishMaster, "圆梦大师", "作为圆梦人完成 20 个心愿", "👑"); err != nil {
			return err
		}
	}
	return nil
}

func (s *badgeService) ListMine(userID uint64) ([]dto.BadgeResponse, error) {
	badges, err := s.badge.ListByUserID(userID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items := make([]dto.BadgeResponse, 0, len(badges))
	for i := range badges {
		items = append(items, dto.ToBadgeResponse(&badges[i]))
	}
	return items, nil
}

func (s *badgeService) Leaderboard(limit int) ([]model.FulfillerStat, error) {
	cacheKey := fmt.Sprintf("wishwall:leaderboard:%d", limit)
	ctx := context.Background()
	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
			var rows []model.FulfillerStat
			if json.Unmarshal([]byte(cached), &rows) == nil {
				s.logger.Info(constants.LogLeaderboardBuilt, "source", "redis", "count", len(rows))
				return rows, nil
			}
		}
	}
	rows, err := s.wish.TopFulfillers(limit)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if s.redis != nil {
		if b, merr := json.Marshal(rows); merr == nil {
			_ = s.redis.Set(ctx, cacheKey, b, 60*time.Second).Err()
		}
	}
	s.logger.Info(constants.LogLeaderboardBuilt, "source", "db", "count", len(rows))
	return rows, nil
}
