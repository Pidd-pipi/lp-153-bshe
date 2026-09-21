package service

import (
	"errors"
	"log/slog"

	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/dto"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/repository"
	"github.com/wishwall/wishwall/internal/util"
)

// BlessingService 祝福留言服务：留言板祝福与虚拟礼物。
type BlessingService interface {
	Create(userID, wishID uint64, req dto.CreateBlessingRequest, ip, requestID string) (*model.Blessing, error)
	ListByWish(wishID uint64, q dto.PageQuery) (*dto.PageResult, error)
}

type blessingService struct {
	bless  repository.BlessingRepository
	wish   repository.WishRepository
	user   repository.UserRepository
	badge  BadgeService
	audit  AuditService
	logger *slog.Logger
}

// NewBlessingService 构造祝福服务。
func NewBlessingService(
	bless repository.BlessingRepository,
	wish repository.WishRepository,
	user repository.UserRepository,
	badge BadgeService,
	audit AuditService,
	logger *slog.Logger,
) BlessingService {
	return &blessingService{bless: bless, wish: wish, user: user, badge: badge, audit: audit, logger: logger}
}

func (s *blessingService) Create(userID, wishID uint64, req dto.CreateBlessingRequest, ip, requestID string) (*model.Blessing, error) {
	wish, err := s.wish.FindByID(wishID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	blessing := &model.Blessing{
		WishID:        wishID,
		UserID:        userID,
		Content:       req.Content,
		GiftEmoji:     req.GiftEmoji,
		IsCelebrating: wish.Status == constants.WishStatusCompleted,
	}
	if err := s.bless.Create(blessing); err != nil {
		return nil, util.NewAppError(constants.CodeBlessingFailed, "祝福发送失败", err)
	}
	s.logger.Info(constants.LogBlessingCreated, "wish_id", wishID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "create_blessing", EntityType: "blessing", EntityID: u64str(blessing.ID),
		Detail: "送出祝福", IP: ip, RequestID: requestID,
	})
	if err := s.badge.GrantFirstBlessing(userID); err != nil {
		s.logger.Warn("grant first blessing badge failed", "error", err)
	}
	return blessing, nil
}

func (s *blessingService) ListByWish(wishID uint64, q dto.PageQuery) (*dto.PageResult, error) {
	page, size := q.Normalize()
	total, err := s.bless.CountByWishID(wishID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items, err := s.bless.ListByWishID(wishID, (page-1)*size, size)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	resps := make([]dto.BlessingResponse, 0, len(items))
	for i := range items {
		name := ""
		if u, uerr := s.user.FindByID(items[i].UserID); uerr == nil {
			name = u.Nickname
		}
		resps = append(resps, dto.ToBlessingResponse(&items[i], name))
	}
	return &dto.PageResult{Items: resps, Total: total, Page: page, PageSize: size}, nil
}
