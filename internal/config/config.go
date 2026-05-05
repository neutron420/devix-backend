package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Media    MediaConfig
	R2       R2Config
	CORS     CORSConfig
	Rate     RateLimitConfig
}

type ServerConfig struct {
	Port         string
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	URL      string
	MaxConns int32
	MinConns int32
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type MediaConfig struct {
	UploadDir         string
	MaxImageSize      int64
	MaxVideoSize      int64
	AllowedImageTypes []string
	AllowedVideoTypes []string
	MaxImagesPerPost  int
	StorageType       string
}

type R2Config struct {
	AccountID  string
	AccessKey  string
	SecretKey  string
	BucketName string
	PublicURL  string
	Endpoint   string
}

type CORSConfig struct {
	Origins []string
}

type RateLimitConfig struct {
	Requests     int
	Window       time.Duration
	AuthRequests int
	AuthWindow   time.Duration
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Env:          getEnv("SERVER_ENV", "development"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			URL:      mustGetEnv("DATABASE_URL"),
			MaxConns: int32(getIntEnv("DATABASE_MAX_CONNS", 25)),
			MinConns: int32(getIntEnv("DATABASE_MIN_CONNS", 5)),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:  mustGetEnv("JWT_ACCESS_SECRET"),
			RefreshSecret: mustGetEnv("JWT_REFRESH_SECRET"),
			AccessExpiry:  getDurationEnv("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 168*time.Hour),
		},
		Media: MediaConfig{
			UploadDir:         getEnv("UPLOAD_DIR", "./uploads"),
			MaxImageSize:      getInt64Env("MAX_IMAGE_SIZE", 5<<20),
			MaxVideoSize:      getInt64Env("MAX_VIDEO_SIZE", 50<<20),
			AllowedImageTypes: getSliceEnv("ALLOWED_IMAGE_TYPES", []string{"image/jpeg", "image/png", "image/webp"}),
			AllowedVideoTypes: getSliceEnv("ALLOWED_VIDEO_TYPES", []string{"video/mp4", "video/webm"}),
			MaxImagesPerPost:  getIntEnv("MAX_IMAGES_PER_POST", 10),
			StorageType:       getEnv("STORAGE_TYPE", "local"),
		},
		R2: R2Config{
			AccountID:  getEnv("R2_ACCOUNT_ID", ""),
			AccessKey:  getEnv("R2_ACCESS_KEY", ""),
			SecretKey:  getEnv("R2_SECRET_KEY", ""),
			BucketName: getEnv("R2_BUCKET_NAME", ""),
			PublicURL:  getEnv("R2_PUBLIC_URL", ""),
			Endpoint:   getEnv("R2_ENDPOINT", ""),
		},
		CORS: CORSConfig{
			Origins: getSliceEnv("CORS_ORIGINS", []string{"http://localhost:3000"}),
		},
		Rate: RateLimitConfig{
			Requests:     getIntEnv("RATE_LIMIT_REQUESTS", 100),
			Window:       getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),
			AuthRequests: getIntEnv("AUTH_RATE_LIMIT_REQUESTS", 5),
			AuthWindow:   getDurationEnv("AUTH_RATE_LIMIT_WINDOW", time.Minute),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) IsProd() bool {
	return c.Server.Env == "production"
}

func (c *Config) IsDev() bool {
	return c.Server.Env == "development"
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWT.AccessSecret == "" || c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT secrets are required")
	}
	if len(c.JWT.AccessSecret) < 32 {
		return fmt.Errorf("JWT_ACCESS_SECRET must be at least 32 characters")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return val
}

func getIntEnv(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}

func getInt64Env(key string, fallback int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}
	return i
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}

func getSliceEnv(key string, fallback []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
