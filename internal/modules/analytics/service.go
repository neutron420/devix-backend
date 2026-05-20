package analytics

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"devix-backend/internal/models"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"github.com/rs/zerolog"
)

type Service struct {
	repo *Repository
	log  zerolog.Logger
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With().Str("module", "analytics").Logger(),
	}
}

func (s *Service) TrackView(ctx context.Context, targetID uuid.UUID, uaString, ip, country, referrer string) {
	ua := user_agent.New(uaString)
	browser, version := ua.Browser()
	ipHash := ""
	if strings.TrimSpace(ip) != "" {
		sum := sha256.Sum256([]byte(ip))
		ipHash = hex.EncodeToString(sum[:])
	}

	event := &models.AnalyticsEvent{
		ID:        uuid.New(),
		TargetID:  targetID,
		Type:      "view",
		Country:   strings.TrimSpace(country),
		Browser:   browser + " " + version,
		OS:        ua.OS(),
		Device:    "Desktop",
		Referrer:  referrer,
		IPHash:    ipHash,
		CreatedAt: time.Now(),
	}
	if ua.Mobile() {
		event.Device = "Mobile"
	} else if ua.Bot() {
		event.Device = "Bot"
	}

	_ = s.repo.LogEvent(ctx, event)
}

func (s *Service) GetPostStats(ctx context.Context, postID uuid.UUID) (*PostAnalyticsResponse, error) {
	return s.repo.GetPostStats(ctx, postID)
}
