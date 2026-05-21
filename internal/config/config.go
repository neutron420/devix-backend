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
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Elasticsearch ElasticsearchConfig
	JWT           JWTConfig
	Media         MediaConfig
	R2            R2Config
	CORS          CORSConfig
	Rate          RateLimitConfig
	Password      PasswordPolicyConfig
	AuthLockout   AuthLockoutConfig
	Pagination    PaginationConfig
	Metrics       MetricsConfig
	Email         EmailConfig
}

type EmailConfig struct {
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	FrontendURL string
}

type ElasticsearchConfig struct {
	URL      string
	Username string
	Password string
}

type ServerConfig struct {
	Port         string
	Env          string
	LogLevel     string
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
	UploadDir           string
	MaxImageSize        int64
	MaxVideoSize        int64
	MaxDocSize          int64
	AllowedImageTypes   []string
	AllowedVideoTypes   []string
	AllowedDocTypes     []string
	MaxImagesPerPost    int
	StorageType         string
	ImageTransformQuery string
}

type R2Config struct {
	AccountID  string
	AccessKey  string
	SecretKey  string
	BucketName string
	PublicURL  string
	CDNURL     string
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

type PasswordPolicyConfig struct {
	MinLength     int
	RequireUpper  bool
	RequireLower  bool
	RequireNumber bool
	RequireSymbol bool
}

type AuthLockoutConfig struct {
	MaxAttempts int
	Window      time.Duration
	BaseLock    time.Duration
	MaxLock     time.Duration
}

type PaginationConfig struct {
	CursorSecret string
}

type MetricsConfig struct {
	Enabled bool
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Env:          getEnv("SERVER_ENV", "development"),
			LogLevel:     getEnv("LOG_LEVEL", "info"),
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
			UploadDir:           getEnv("UPLOAD_DIR", "./uploads"),
			MaxImageSize:        getInt64Env("MAX_IMAGE_SIZE", 5<<20),
			MaxVideoSize:        getInt64Env("MAX_VIDEO_SIZE", 50<<20),
			MaxDocSize:          getInt64Env("MAX_DOC_SIZE", 10<<20),
			AllowedImageTypes:   getSliceEnv("ALLOWED_IMAGE_TYPES", []string{"image/jpeg", "image/png", "image/webp"}),
			AllowedVideoTypes:   getSliceEnv("ALLOWED_VIDEO_TYPES", []string{"video/mp4", "video/webm"}),
			AllowedDocTypes:     getSliceEnv("ALLOWED_DOC_TYPES", []string{"application/pdf"}),
			MaxImagesPerPost:    getIntEnv("MAX_IMAGES_PER_POST", 10),
			StorageType:         getEnv("STORAGE_TYPE", "local"),
			ImageTransformQuery: getEnv("IMAGE_CDN_QUERY", ""),
		},
		R2: R2Config{
			AccountID:  getEnv("R2_ACCOUNT_ID", ""),
			AccessKey:  getEnv("R2_ACCESS_KEY", ""),
			SecretKey:  getEnv("R2_SECRET_KEY", ""),
			BucketName: getEnv("R2_BUCKET_NAME", ""),
			PublicURL:  getEnv("R2_PUBLIC_URL", ""),
			CDNURL:     getEnv("R2_CDN_URL", ""),
			Endpoint:   getEnv("R2_ENDPOINT", ""),
		},
		Elasticsearch: ElasticsearchConfig{
			URL: getEnv("ELASTICSEARCH_URL", ""),
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
		Password: PasswordPolicyConfig{
			MinLength:     getIntEnv("PASSWORD_MIN_LENGTH", 10),
			RequireUpper:  getBoolEnv("PASSWORD_REQUIRE_UPPER", true),
			RequireLower:  getBoolEnv("PASSWORD_REQUIRE_LOWER", true),
			RequireNumber: getBoolEnv("PASSWORD_REQUIRE_NUMBER", true),
			RequireSymbol: getBoolEnv("PASSWORD_REQUIRE_SYMBOL", true),
		},
		AuthLockout: AuthLockoutConfig{
			MaxAttempts: getIntEnv("AUTH_LOCKOUT_MAX_ATTEMPTS", 5),
			Window:      getDurationEnv("AUTH_LOCKOUT_WINDOW", 15*time.Minute),
			BaseLock:    getDurationEnv("AUTH_LOCKOUT_BASE", 2*time.Minute),
			MaxLock:     getDurationEnv("AUTH_LOCKOUT_MAX", 30*time.Minute),
		},
		Pagination: PaginationConfig{
			CursorSecret: getEnv("CURSOR_SECRET", ""),
		},
		Metrics: MetricsConfig{
			Enabled: getBoolEnv("METRICS_ENABLED", true),
		},
		Email: EmailConfig{
			SMTPHost:    getEnv("SMTP_HOST", ""),
			SMTPPort:    getIntEnv("SMTP_PORT", 587),
			SMTPUser:    getEnv("SMTP_USER", ""),
			SMTPPass:    getEnv("SMTP_PASS", ""),
			SMTPFrom:    getEnv("SMTP_FROM", "noreply@devix.app"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
	}

	if cfg.Pagination.CursorSecret == "" {
		cfg.Pagination.CursorSecret = cfg.JWT.AccessSecret
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

func getBoolEnv(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return b
}
