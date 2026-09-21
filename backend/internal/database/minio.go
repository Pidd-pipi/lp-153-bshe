package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/wishwall/wishwall/internal/config"
	"github.com/wishwall/wishwall/internal/constants"
)

// OpenMinIO 建立 MinIO 客户端，并确保 bucket 存在。
func OpenMinIO(cfg *config.Config, logger *slog.Logger) (*minio.Client, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		logger.Error(constants.LogMinIOInitFailed, "endpoint", cfg.MinIOEndpoint, "error", err)
		return nil, fmt.Errorf("init minio client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		return nil, fmt.Errorf("check minio bucket %s: %w", cfg.MinIOBucket, err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.MinIOBucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make minio bucket %s: %w", cfg.MinIOBucket, err)
		}
	}
	// 开放桶匿名只读，使上传的图片/音频可通过 URL 直接访问。
	publicPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + cfg.MinIOBucket + `/*"]}]}`
	if err := client.SetBucketPolicy(ctx, cfg.MinIOBucket, publicPolicy); err != nil {
		return nil, fmt.Errorf("set minio bucket policy %s: %w", cfg.MinIOBucket, err)
	}
	logger.Info("minio client ready", "bucket", cfg.MinIOBucket)
	return client, nil
}
