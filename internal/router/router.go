package router

import (
	"devix-backend/internal/config"
	"devix-backend/internal/middleware"
	"devix-backend/internal/modules/auth"
	"devix-backend/internal/modules/comment"
	"devix-backend/internal/modules/media"
	"devix-backend/internal/modules/post"
	"devix-backend/internal/modules/tag"
	"devix-backend/internal/modules/user"
	"devix-backend/internal/modules/vote"
	"devix-backend/internal/modules/websocket"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Handlers holds all module handlers for dependency injection.
type Handlers struct {
	Auth    *auth.Handler
	User    *user.Handler
	Post    *post.Handler
	Comment *comment.Handler
	Tag     *tag.Handler
	Vote    *vote.Handler
	WS      *websocket.Handler
}

// Setup configures the Gin engine with all routes and middleware.
func Setup(cfg *config.Config, log zerolog.Logger, jwtManager *jwtpkg.Manager, handlers *Handlers) *gin.Engine {
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Recovery middleware
	r.Use(gin.Recovery())

	// Global middleware chain
	rateLimiter := middleware.NewInMemoryRateLimiter()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS(cfg.CORS.Origins))
	r.Use(middleware.RateLimit(rateLimiter, cfg.Rate.Requests, cfg.Rate.Window))

	// Serve uploaded files
	media.RegisterRoutes(r, cfg.Media.UploadDir)

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Devix API",
			"status":  "operational",
			"version": "1.0.0",
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "devix-backend"})
	})

	// WebSocket route
	r.GET("/ws", middleware.Auth(jwtManager), handlers.WS.ServeWS)

	// API v1 routes
	v1 := r.Group("/api/v1")

	// Auth routes (with stricter rate limiting)
	authGroup := v1.Group("")
	authGroup.Use(middleware.AuthRateLimit(rateLimiter, cfg.Rate.AuthRequests, cfg.Rate.AuthWindow))
	auth.RegisterRoutes(authGroup, handlers.Auth)

	// Module routes
	user.RegisterRoutes(v1, handlers.User, jwtManager)
	post.RegisterRoutes(v1, handlers.Post, jwtManager)
	comment.RegisterRoutes(v1, handlers.Comment, jwtManager)
	tag.RegisterRoutes(v1, handlers.Tag)
	vote.RegisterRoutes(v1, handlers.Vote, jwtManager)

	return r
}
