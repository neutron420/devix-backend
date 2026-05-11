package analytics

import (
	"context"
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

func (s *Service) TrackView(ctx context.Context, targetID uuid.UUID, uaString, ip, referrer string) {
	ua := user_agent.New(uaString)
	browser, version := ua.Browser()
	
	event := &models.AnalyticsEvent{
		ID:        uuid.New(),
		TargetID:  targetID,
		Type:      "view",
		Browser:   browser + " " + version,
		OS:        ua.OS(),
		Device:    "Desktop", // Simplified
		Referrer:  referrer,
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
