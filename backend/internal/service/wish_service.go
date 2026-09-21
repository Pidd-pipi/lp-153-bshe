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

// WishService 心愿服务：发布/更新/删除/查询/点赞/发现广场。
type WishService interface {
	Create(userID uint64, req dto.CreateWishRequest, ip, requestID string) (*model.Wish, error)
	Update(userID, wishID uint64, req dto.UpdateWishRequest, ip, requestID string) (*model.Wish, error)
	Delete(userID, wishID uint64, ip, requestID string) error
	GetByID(wishID uint64) (*dto.WishDetailResponse, error)
	List(q dto.WishQuery) (*dto.PageResult, error)
	ListMine(userID uint64, q dto.PageQuery) (*dto.PageResult, error)
	Like(userID, wishID uint64, ip, requestID string) error
	Discover(storyQuery dto.PageQuery) (*dto.PageResult, error)
	Leaderboard(limit int) ([]model.FulfillerStat, error)
}

type wishService struct {
	wish   repository.WishRepository
	claim  repository.WishClaimRepository
	bless  repository.BlessingRepository
	user   repository.UserRepository
	badge  BadgeService
	audit  AuditService
	logger *slog.Logger
}

// NewWishService 构造心愿服务。
func NewWishService(
	wish repository.WishRepository,
	claim repository.WishClaimRepository,
	bless repository.BlessingRepository,
	user repository.UserRepository,
	badge BadgeService,
	audit AuditService,
	logger *slog.Logger,
) WishService {
	return &wishService{wish: wish, claim: claim, bless: bless, user: user, badge: badge, audit: audit, logger: logger}
}

func (s *wishService) Create(userID uint64, req dto.CreateWishRequest, ip, requestID string) (*model.Wish, error) {
	wish := &model.Wish{
		UserID:           userID,
		Title:            req.Title,
		Content:          req.Content,
		ImageURLs:        util.EncodeStringArray(req.ImageURLs),
		Category:         req.Category,
		Visibility:       req.Visibility,
		Difficulty:       req.Difficulty,
		ExpectedDeadline: req.ExpectedDeadline,
		Status:           constants.WishStatusPending,
		IsAnonymous:      req.IsAnonymous,
	}
	if err := s.wish.Create(wish); err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogWishCreated, "wish_id", wish.ID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "create_wish", EntityType: "wish", EntityID: u64str(wish.ID),
		Detail: "发布心愿：" + util.TruncateString(wish.Title, 30), IP: ip, RequestID: requestID,
	})
	if err := s.badge.GrantFirstWish(userID); err != nil {
		s.logger.Warn("grant first wish badge failed", "error", err)
	}
	return wish, nil
}

func (s *wishService) Update(userID, wishID uint64, req dto.UpdateWishRequest, ip, requestID string) (*model.Wish, error) {
	wish, err := s.wish.FindByID(wishID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if wish.UserID != userID {
		return nil, util.NewAppError(constants.CodeWishNotOwner, constants.MsgWishNotOwner, errors.New("wish owner mismatch"))
	}
	if req.Title != "" {
		wish.Title = req.Title
	}
	if req.Content != "" {
		wish.Content = req.Content
	}
	if req.ImageURLs != nil {
		wish.ImageURLs = util.EncodeStringArray(req.ImageURLs)
	}
	if req.Category != "" {
		wish.Category = req.Category
	}
	if req.Visibility != "" {
		wish.Visibility = req.Visibility
	}
	if req.Difficulty != "" {
		wish.Difficulty = req.Difficulty
	}
	if req.ExpectedDeadline != nil {
		wish.ExpectedDeadline = req.ExpectedDeadline
	}
	if req.IsAnonymous != nil {
		wish.IsAnonymous = *req.IsAnonymous
	}
	if err := s.wish.Update(wish); err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogWishUpdated, "wish_id", wish.ID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "update_wish", EntityType: "wish", EntityID: u64str(wish.ID),
		Detail: "更新心愿", IP: ip, RequestID: requestID,
	})
	return wish, nil
}

func (s *wishService) Delete(userID, wishID uint64, ip, requestID string) error {
	wish, err := s.wish.FindByID(wishID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if wish.UserID != userID {
		return util.NewAppError(constants.CodeWishNotOwner, constants.MsgWishNotOwner, errors.New("wish owner mismatch"))
	}
	if err := s.wish.Delete(wishID); err != nil {
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogWishDeleted, "wish_id", wish.ID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "delete_wish", EntityType: "wish", EntityID: u64str(wishID),
		Detail: "删除心愿", IP: ip, RequestID: requestID,
	})
	return nil
}

func (s *wishService) GetByID(wishID uint64) (*dto.WishDetailResponse, error) {
	wish, err := s.wish.FindByID(wishID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	detail := &dto.WishDetailResponse{WishResponse: dto.ToWishResponse(wish)}
	if author, aerr := s.user.FindByID(wish.UserID); aerr == nil {
		detail.AuthorNickname = author.Nickname
		detail.AuthorAvatar = author.Avatar
	}
	if claim, cerr := s.claim.FindByWishID(wishID); cerr == nil {
		claimResp := dto.ToWishClaimResponse(claim, wish.Title, "")
		if fulfiller, ferr := s.user.FindByID(claim.UserID); ferr == nil {
			claimResp.FulfillerName = fulfiller.Nickname
		}
		detail.Claim = &claimResp
	}
	if count, berr := s.bless.CountByWishID(wishID); berr == nil {
		detail.BlessingCount = count
	}
	return detail, nil
}

func (s *wishService) List(q dto.WishQuery) (*dto.PageResult, error) {
	page, size := q.Normalize()
	filters := map[string]any{
		"keyword":    q.Keyword,
		"category":   q.Category,
		"status":     q.Status,
		"visibility": q.Visibility,
		"user_id":    q.UserID,
	}
	total, err := s.wish.Count(filters)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	wishes, err := s.wish.List(filters, q.Sort, (page-1)*size, size)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items, err := s.assembleWishes(wishes)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: size}, nil
}

func (s *wishService) ListMine(userID uint64, q dto.PageQuery) (*dto.PageResult, error) {
	page, size := q.Normalize()
	total, err := s.wish.CountByUserID(userID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	wishes, err := s.wish.ListByUserID(userID, (page-1)*size, size)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items, err := s.assembleWishes(wishes)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: size}, nil
}

// assembleWishes 批量组装心愿返回（作者信息 + 匿名隐藏）。
func (s *wishService) assembleWishes(wishes []model.Wish) ([]dto.WishResponse, error) {
	items := make([]dto.WishResponse, 0, len(wishes))
	for i := range wishes {
		resp := dto.ToWishResponse(&wishes[i])
		if wishes[i].IsAnonymous {
			resp.AuthorNickname = "匿名心愿"
		} else if author, aerr := s.user.FindByID(wishes[i].UserID); aerr == nil {
			resp.AuthorNickname = author.Nickname
			resp.AuthorAvatar = author.Avatar
		}
		items = append(items, resp)
	}
	return items, nil
}

func (s *wishService) Like(userID, wishID uint64, ip, requestID string) error {
	if _, err := s.wish.FindByID(wishID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeWishNotFound, constants.MsgWishNotFound, err)
		}
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	if err := s.wish.IncrementLikes(wishID); err != nil {
		return util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogWishLike, "wish_id", wishID, "user_id", userID)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "like_wish", EntityType: "wish", EntityID: u64str(wishID),
		Detail: "点赞心愿", IP: ip, RequestID: requestID,
	})
	return nil
}

// Discover 发现广场：最新完成的心愿故事（复用 List 的 status=completed 过滤 + BadgeService.Leaderboard）。
func (s *wishService) Discover(storyQuery dto.PageQuery) (*dto.PageResult, error) {
	page, size := storyQuery.Normalize()
	filters := map[string]any{"status": constants.WishStatusCompleted}
	total, err := s.wish.Count(filters)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	stories, err := s.wish.List(filters, "new", (page-1)*size, size)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	items, err := s.assembleWishes(stories)
	if err != nil {
		return nil, util.NewAppError(constants.CodeInternalError, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogLeaderboardBuilt, "stories", len(items))
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: size}, nil
}

// Leaderboard 圆梦人排行榜（复用 BadgeService.Leaderboard）。
func (s *wishService) Leaderboard(limit int) ([]model.FulfillerStat, error) {
	return s.badge.Leaderboard(limit)
}
