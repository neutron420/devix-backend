package router

import (
	"time"

	"devix-backend/internal/config"
	"devix-backend/internal/middleware"
	"devix-backend/internal/modules/audit"
	"devix-backend/internal/modules/auth"
	"devix-backend/internal/modules/bookmark"
	"devix-backend/internal/modules/comment"
	"devix-backend/internal/modules/follow"
	"devix-backend/internal/modules/media"
	"devix-backend/internal/modules/notification"
	"devix-backend/internal/modules/post"
	"devix-backend/internal/modules/tag"
	"devix-backend/internal/modules/user"
	"devix-backend/internal/modules/vote"
	"devix-backend/internal/modules/websocket"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Handlers struct {
	Auth    *auth.Handler
	User    *user.Handler
	Post    *post.Handler
	Comment *comment.Handler
	Tag          *tag.Handler
	Vote         *vote.Handler
	Notification *notification.Handler
	Bookmark     *bookmark.Handler
	Follow       *follow.Handler
	WS           *websocket.Handler
}

func Setup(cfg *config.Config, log zerolog.Logger, jwtManager *jwtpkg.Manager, handlers *Handlers, redisClient *redis.Client, auditSvc *audit.Service) *gin.Engine {
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())

	// Security Headers (CSP, HSTS, X-Frame-Options, etc.)
	r.Use(middleware.SecurityHeaders())

	// Abuse Detection
	abuseDetector := middleware.NewAbuseDetector(redisClient, 50, 10*time.Minute)
	r.Use(middleware.AbuseProtection(abuseDetector))

	// Redis-backed Rate Limiting (falls back to in-memory if Redis is nil)
	redisLimiter := middleware.NewRedisRateLimiter(redisClient)

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS(cfg.CORS.Origins))
	r.Use(middleware.RedisRateLimit(redisLimiter, cfg.Rate.Requests, cfg.Rate.Window))

	media.RegisterRoutes(r, cfg.Media.UploadDir)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Devix API",
			"status":  "operational",
			"version": "1.0.0",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "devix-backend"})
	})

	r.GET("/ws", middleware.Auth(jwtManager), handlers.WS.ServeWS)

	v1 := r.Group("/api/v1")

	authGroup := v1.Group("")
	authGroup.Use(middleware.AuthRateLimitRedis(redisLimiter, cfg.Rate.AuthRequests, cfg.Rate.AuthWindow))
	authGroup.Use(middleware.LoginProtection(redisClient))
	auth.RegisterRoutes(authGroup, handlers.Auth)

	user.RegisterRoutes(v1, handlers.User, jwtManager)
	post.RegisterRoutes(v1, handlers.Post, jwtManager)
	comment.RegisterRoutes(v1, handlers.Comment, jwtManager)
	tag.RegisterRoutes(v1, handlers.Tag)
	vote.RegisterRoutes(v1, handlers.Vote, jwtManager)
	notification.RegisterRoutes(v1, handlers.Notification, jwtManager)
	bookmark.RegisterRoutes(v1, handlers.Bookmark, jwtManager)
	follow.RegisterRoutes(v1, handlers.Follow, jwtManager)

	_ = auditSvc

	return r
}
