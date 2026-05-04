package tag

import (
	"context"
	"strings"

	apperrors "devix-backend/internal/errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo *Repository
	log  zerolog.Logger
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{repo: repo, log: log.With().Str("module", "tag").Logger()}
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
