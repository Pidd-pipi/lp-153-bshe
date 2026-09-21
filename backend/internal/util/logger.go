package util

import (
	"log/slog"
	"os"
	"sync"
)

var (
	loggerOnce sync.Once
	appLogger  *slog.Logger
)

// InitLogger 初始化全局 slog 日志器。JSON 结构化输出，development 环境开启 Debug 级别。
func InitLogger(env string) *slog.Logger {
	loggerOnce.Do(func() {
		level := slog.LevelInfo
		if env == "development" {
			level = slog.LevelDebug
		}
		appLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	})
	return appLogger
}

// Logger 返回全局日志器，未初始化时按 development 兜底初始化。
func Logger() *slog.Logger {
	if appLogger == nil {
		return InitLogger("development")
	}
	return appLogger
}
