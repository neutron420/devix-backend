package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devix-backend/internal/config"
	"devix-backend/internal/database"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/auth"
	"devix-backend/internal/modules/comment"
	"devix-backend/internal/modules/media"
	"devix-backend/internal/modules/post"
	"devix-backend/internal/modules/tag"
	"devix-backend/internal/modules/user"
	"devix-backend/internal/modules/vote"
	wsmod "devix-backend/internal/modules/websocket"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"devix-backend/internal/pkg/logger"
	"devix-backend/internal/queue"
	"devix-backend/internal/router"
	"devix-backend/internal/validator"
)

func main() {
	// ── 1. Load configuration ──────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// ── 2. Initialize logger ───────────────────────────────────────
	log := logger.New(cfg.Server.Env)
	log.Info().
		Str("env", cfg.Server.Env).
		Str("port", cfg.Server.Port).
		Msg("starting devix-backend")

	ctx := context.Background()

	// ── 3. Connect to PostgreSQL (GORM) ───────────────────────────
	db, err := database.NewGormDB(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to PostgreSQL")
	}
	log.Info().Msg("connected to PostgreSQL (GORM)")

	// ── 3b. Auto Migration ────────────────────────────────────────
	log.Info().Msg("running auto-migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Media{},
		&models.Comment{},
		&models.Tag{},
		&models.Vote{},
		&models.RefreshToken{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("auto-migration failed")
	}
	log.Info().Msg("auto-migration completed successfully")

	// ── 4. Connect to Redis (optional) ─────────────────────────────
	redis, err := database.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Warn().Err(err).Msg("Redis unavailable — running without cache")
	} else if redis != nil {
		defer redis.Close()
		log.Info().Msg("connected to Redis")
	} else {
		log.Info().Msg("Redis not configured — running without cache")
	}

	// ── 5. Setup custom validators ─────────────────────────────────
	validator.Setup()

	// ── 6. Initialize JWT Manager ──────────────────────────────────
	jwtManager := jwtpkg.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	// ── 7. Initialize storage provider ─────────────────────────────
	var storage media.StorageProvider
	if cfg.Media.StorageType == "r2" {
		var err error
		storage, err = media.NewR2Storage(ctx, cfg.R2)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize R2 storage")
		}
		log.Info().Msg("using Cloudflare R2 storage")
	} else {
		storage = media.NewLocalStorage(cfg.Media.UploadDir, "/uploads")
		log.Info().Msg("using local filesystem storage")
	}

	// ── 8. Initialize background job queue ─────────────────────────
	jobQueue := queue.New(100, log)
	jobQueue.Start(ctx, 3)
	defer jobQueue.Close()

	// ── 8b. Initialize WebSocket Hub ──────────────────────────────
	hub := wsmod.NewHub()
	go hub.Run()

	// ── 9. Initialize repositories ─────────────────────────────────
	authRepo := auth.NewRepository(db)
	userRepo := user.NewRepository(db)
	mediaRepo := media.NewRepository(db)
	postRepo := post.NewRepository(db)
	commentRepo := comment.NewRepository(db)
	tagRepo := tag.NewRepository(db)
	voteRepo := vote.NewRepository(db)

	// ── 10. Initialize services ────────────────────────────────────
	wsService := wsmod.NewService(hub, log)
	authService := auth.NewService(authRepo, jwtManager, log)
	mediaService := media.NewService(mediaRepo, storage, cfg.Media, log)
	userService := user.NewService(userRepo, log)
	tagService := tag.NewService(tagRepo, log)
	postService := post.NewService(postRepo, mediaService, tagService, log)
	commentService := comment.NewService(commentRepo, wsService, log)
	voteService := vote.NewService(voteRepo, log)

	// ── 11. Initialize handlers ────────────────────────────────────
	handlers := &router.Handlers{
		Auth:    auth.NewHandler(authService),
		User:    user.NewHandler(userService, mediaService),
		Post:    post.NewHandler(postService, mediaService),
		Comment: comment.NewHandler(commentService),
		Tag:     tag.NewHandler(tagService),
		Vote:    vote.NewHandler(voteService),
		WS:      wsmod.NewHandler(hub),
	}

	// ── 12. Setup router ───────────────────────────────────────────
	engine := router.Setup(cfg, log, jwtManager, handlers)

	// ── 13. Start HTTP server ──────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	log.Info().Str("addr", srv.Addr).Msg("server started")

	// ── 14. Wait for shutdown signal ───────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Info().Str("signal", sig.String()).Msg("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server stopped gracefully")
}
