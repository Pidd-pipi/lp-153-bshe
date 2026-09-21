package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/wishwall/wishwall/internal/config"
	"github.com/wishwall/wishwall/internal/constants"
	"github.com/wishwall/wishwall/internal/model"
	"github.com/wishwall/wishwall/internal/util"
)

var allowedImageExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
var allowedAudioExt = map[string]bool{".mp3": true, ".wav": true, ".m4a": true, ".aac": true}

// UploadService 文件上传服务：图片/音频写入 MinIO 并返回可访问 URL。
type UploadService interface {
	Upload(ctx context.Context, userID uint64, file multipart.File, header *multipart.FileHeader, kind, ip, requestID string) (string, error)
}

type uploadService struct {
	cfg    *config.Config
	minio  *minio.Client
	audit  AuditService
	logger *slog.Logger
}

// NewUploadService 构造上传服务。
func NewUploadService(cfg *config.Config, client *minio.Client, audit AuditService, logger *slog.Logger) UploadService {
	return &uploadService{cfg: cfg, minio: client, audit: audit, logger: logger}
}

func (s *uploadService) Upload(ctx context.Context, userID uint64, file multipart.File, header *multipart.FileHeader, kind, ip, requestID string) (string, error) {
	if header.Size > 10<<20 {
		return "", util.NewAppError(constants.CodeUploadFailed, "文件大小不能超过 10MB", fmt.Errorf("file too large"))
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch kind {
	case "image":
		if !allowedImageExt[ext] {
			s.logger.Warn(constants.LogUploadDenied, "kind", kind, "ext", ext)
			return "", util.NewAppError(constants.CodeFileTypeNotAllowed, "仅支持 jpg/png/webp/gif 图片", fmt.Errorf("image ext %s not allowed", ext))
		}
	case "audio":
		if !allowedAudioExt[ext] {
			s.logger.Warn(constants.LogUploadDenied, "kind", kind, "ext", ext)
			return "", util.NewAppError(constants.CodeFileTypeNotAllowed, "仅支持 mp3/wav/m4a/aac 音频", fmt.Errorf("audio ext %s not allowed", ext))
		}
	default:
		return "", util.NewAppError(constants.CodeBadRequest, "上传类型必须为 image 或 audio", fmt.Errorf("unknown upload kind %s", kind))
	}
	objectKey := fmt.Sprintf("uploads/%s/%s/%s%s", kind, uuid.NewString(), uuid.NewString(), ext)
	contentType := "application/octet-stream"
	if kind == "image" {
		contentType = "image/" + strings.TrimPrefix(ext, ".")
	} else if kind == "audio" {
		contentType = "audio/" + strings.TrimPrefix(ext, ".")
	}
	_, err := s.minio.PutObject(ctx, s.cfg.MinIOBucket, objectKey, file, header.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", util.NewAppError(constants.CodeUploadFailed, "文件上传失败", err)
	}
	url := util.BuildUploadURL(s.cfg.UploadBaseURL, s.cfg.MinIOBucket, objectKey)
	s.logger.Info(constants.LogFileUploaded, "user_id", userID, "object", objectKey)
	_ = s.audit.Record(&model.AuditLog{
		UserID: userID, Action: "upload_file", EntityType: "upload", EntityID: objectKey,
		Detail: "上传文件 kind=" + kind, IP: ip, RequestID: requestID,
	})
	return url, nil
}

