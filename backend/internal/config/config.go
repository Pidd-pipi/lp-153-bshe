package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

// Config 集中承载后端全部可配置项，全部通过环境变量注入。
type Config struct {
	AppName    string `env:"APP_NAME" envDefault:"wishwall"`
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`

	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"wishwall_user"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"wishwall_pwd"`
	DBName     string `env:"DB_NAME" envDefault:"wishwall_db"`
	DBSSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`

	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`

	JWTSecret        string        `env:"JWT_SECRET" envDefault:"change_me_to_a_long_random_string"`
	JWTExpireHours   int           `env:"JWT_EXPIRE_HOURS" envDefault:"72"`
	RateLimitPerMin  int           `env:"RATE_LIMIT_PER_MIN" envDefault:"120"`
	ShutdownTimeout  time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`

	MinIOEndpoint  string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	MinIOAccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	MinIOSecretKey string `env:"MINIO_SECRET_KEY" envDefault:"minioadmin"`
	MinIOBucket    string `env:"MINIO_BUCKET" envDefault:"wishwall"`
	MinIOUseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`

	UploadBaseURL string `env:"UPLOAD_BASE_URL" envDefault:""`
}

// Load 从环境变量解析配置，字段缺失时使用 envDefault 兜底。
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
