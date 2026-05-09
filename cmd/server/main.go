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
	"devix-backend/internal/modules/audit"
	"devix-backend/internal/modules/auth"
	"devix-backend/internal/modules/bookmark"
	"devix-backend/internal/modules/comment"
	"devix-backend/internal/modules/follow"
	"devix-backend/internal/modules/media"
	"devix-backend/internal/modules/notification"
	"devix-backend/internal/modules/post"
	"devix-backend/internal/modules/search"
	"devix-backend/internal/modules/tag"
	"devix-backend/internal/modules/user"
	"devix-backend/internal/modules/vote"
	wsmod "devix-backend/internal/modules/websocket"
	"devix-backend/internal/pkg/cache"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"devix-backend/internal/pkg/logger"
	"devix-backend/internal/queue"
	"devix-backend/internal/router"
	"devix-backend/internal/validator"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	log := logger.New(cfg.Server.Env)
	log.Info().
		Str("env", cfg.Server.Env).
		Str("port", cfg.Server.Port).
		Msg("starting devix-backend")

	ctx := context.Background()

	db, err := database.NewGormDB(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to PostgreSQL")
	}
	log.Info().Msg("connected to PostgreSQL (GORM)")

	log.Info().Msg("running auto-migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Media{},
		&models.Comment{},
		&models.Tag{},
		&models.Vote{},
		&models.RefreshToken{},
		&models.Notification{},
		&models.Bookmark{},
		&models.Follow{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("auto-migration failed")
	}
	log.Info().Msg("auto-migration completed successfully")

	redis, err := database.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Warn().Err(err).Msg("Redis unavailable — running without cache")
	} else if redis != nil {
		defer redis.Close()
		log.Info().Msg("connected to Redis")
	} else {
		log.Info().Msg("Redis not configured — running without cache")
	}

	esClient, err := search.NewClient(cfg.Elasticsearch.URL, log)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to elasticsearch, search features will be limited")
	}

	appCache := cache.New(redis)

	validator.Setup()

	jwtManager := jwtpkg.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

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

	jobQueue := queue.New(100, log)
	jobQueue.Start(ctx, 3)
	defer jobQueue.Close()

	hub := wsmod.NewHub(redis, log)
	go hub.Run(ctx)

	authRepo := auth.NewRepository(db)
	userRepo := user.NewRepository(db)
	mediaRepo := media.NewRepository(db)
	postRepo := post.NewRepository(db)
	commentRepo := comment.NewRepository(db)
	tagRepo := tag.NewRepository(db)
	voteRepo := vote.NewRepository(db)
	notificationRepo := notification.NewRepository(db)
	bookmarkRepo := bookmark.NewRepository(db)
	followRepo := follow.NewRepository(db)

	wsService := wsmod.NewService(hub, log)
	searchService := search.NewService(esClient, log)
	authService := auth.NewService(authRepo, jwtManager, log)
	mediaService := media.NewService(mediaRepo, storage, cfg.Media, log)
	userService := user.NewService(userRepo, appCache, log)
	tagService := tag.NewService(tagRepo, appCache, log)
	postService := post.NewService(postRepo, mediaService, tagService, searchService, appCache, jobQueue, log)
	notificationService := notification.NewService(notificationRepo, wsService, log)
	followService := follow.NewService(followRepo, userService, notificationService, log)
	postService.SetFollowService(followService)
	bookmarkService := bookmark.NewService(bookmarkRepo, postService, log)
	commentService := comment.NewService(commentRepo, wsService, notificationService, log)
	voteService := vote.NewService(voteRepo, notificationService, log)
	auditService := audit.NewService(db, log)

	handlers := &router.Handlers{
		Auth:    auth.NewHandler(authService),
		User:    user.NewHandler(userService, mediaService),
		Post:    post.NewHandler(postService, mediaService),
		Comment: comment.NewHandler(commentService),
		Tag:     tag.NewHandler(tagService),
		Vote:         vote.NewHandler(voteService),
		Notification: notification.NewHandler(notificationService),
		Bookmark:     bookmark.NewHandler(bookmarkService),
		Follow:       follow.NewHandler(followService),
		WS:           wsmod.NewHandler(hub),
	}

	engine := router.Setup(cfg, log, jwtManager, handlers, redis, auditService)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	log.Info().Str("addr", srv.Addr).Msg("server started")

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
