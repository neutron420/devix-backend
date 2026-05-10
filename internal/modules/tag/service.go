package tag

import (
	"context"
	"strings"

	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
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

func (s *Service) toResponse(t *models.Tag) TagResponse {
	resp := TagResponse{
		ID:          t.ID.String(),
		Name:        t.Name,
		Slug:        t.Slug,
		Description: t.Description,
		Category:    t.Category,
		Synonyms:    t.Synonyms,
		PostCount:   t.PostCount,
	}
	if t.ParentID != nil {
		pid := t.ParentID.String()
		resp.ParentID = &pid
	}
	return resp
}

func (s *Service) GetAll(ctx context.Context) ([]TagResponse, error) {
	tags, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	result := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		result = append(result, s.toResponse(&t))
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
		result = append(result, s.toResponse(&t))
	}

	_ = s.cache.Set(ctx, cacheKey, result, 15*time.Minute)
	return result, nil
}

func (s *Service) GetTagTree(ctx context.Context) ([]TagTreeResponse, error) {
	cacheKey := "tags:tree"
	var cached []TagTreeResponse
	if ok, _ := s.cache.Get(ctx, cacheKey, &cached); ok {
		return cached, nil
	}

	topLevel, err := s.repo.GetTopLevel(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	tree := make([]TagTreeResponse, 0, len(topLevel))
	for _, t := range topLevel {
		node := TagTreeResponse{
			ID:          t.ID.String(),
			Name:        t.Name,
			Slug:        t.Slug,
			Description: t.Description,
			Category:    t.Category,
			PostCount:   t.PostCount,
		}

		children, err := s.repo.GetChildren(ctx, t.ID)
		if err == nil && len(children) > 0 {
			node.Children = make([]TagTreeResponse, 0, len(children))
			for _, c := range children {
				node.Children = append(node.Children, TagTreeResponse{
					ID:          c.ID.String(),
					Name:        c.Name,
					Slug:        c.Slug,
					Description: c.Description,
					Category:    c.Category,
					PostCount:   c.PostCount,
				})
			}
		}
		tree = append(tree, node)
	}

	_ = s.cache.Set(ctx, cacheKey, tree, 30*time.Minute)
	return tree, nil
}

func (s *Service) GetByCategory(ctx context.Context, category string) ([]TagResponse, error) {
	tags, err := s.repo.GetByCategory(ctx, category)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	result := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		result = append(result, s.toResponse(&t))
	}
	return result, nil
}
