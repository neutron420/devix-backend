package tag

import (
	"context"
	"strings"

	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/pkg/cache"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo  *Repository
	cache *cache.Cache
	log   zerolog.Logger
}

func NewService(repo *Repository, cache *cache.Cache, log zerolog.Logger) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
		log:   log.With().Str("module", "tag").Logger(),
	}
}

func (s *Service) GetAll(ctx context.Context) ([]TagResponse, error) {
	tags, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	result := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		result = append(result, TagResponse{ID: t.ID.String(), Name: t.Name, Slug: t.Slug, Description: t.Description, PostCount: t.PostCount})
	}
	return result, nil
}

func (s *Service) FindOrCreateByNames(ctx context.Context, names []string) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		tag, err := s.repo.GetByName(ctx, name)
		if err != nil {
			return nil, err
		}
		if tag == nil {
			tag, err = s.repo.Create(ctx, name, nil)
			if err != nil {
				s.log.Warn().Err(err).Str("tag", name).Msg("failed to create tag")
				continue
			}
		}
		ids = append(ids, tag.ID)
	}
	return ids, nil
}

func (s *Service) GetTrending(ctx context.Context, limit int) ([]TagResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	cacheKey := "tags:trending"
	var cached []TagResponse
	if ok, _ := s.cache.Get(ctx, cacheKey, &cached); ok {
		return cached, nil
	}

	tags, err := s.repo.GetTrending(ctx, limit, 48*time.Hour)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	result := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		result = append(result, TagResponse{
			ID:          t.ID.String(),
			Name:        t.Name,
			Slug:        t.Slug,
			Description: t.Description,
			PostCount:   t.PostCount,
		})
	}

	_ = s.cache.Set(ctx, cacheKey, result, 15*time.Minute)
	return result, nil
}
