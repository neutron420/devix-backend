package bookmark

import (
	"context"
	"devix-backend/internal/modules/post"
	"devix-backend/internal/pkg/pagination"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo        *Repository
	postService *post.Service
	log         zerolog.Logger
}

func NewService(repo *Repository, postService *post.Service, log zerolog.Logger) *Service {
	return &Service{
		repo:        repo,
		postService: postService,
		log:         log.With().Str("module", "bookmark").Logger(),
	}
}

func (s *Service) ToggleBookmark(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	isBookmarked, err := s.repo.IsBookmarked(ctx, userID, postID)
	if err != nil {
		return false, err
	}

	if isBookmarked {
		if err := s.repo.Unbookmark(ctx, userID, postID); err != nil {
			return false, err
		}
		return false, nil
	}

	if err := s.repo.Bookmark(ctx, userID, postID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) ListBookmarks(ctx context.Context, userID uuid.UUID, cursor string, limit int) (*BookmarkListResponse, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}

	posts, hasMore, err := s.repo.List(ctx, userID, cursor, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]post.PostResponse, 0, len(posts))
	for _, p := range posts {
		// Get detailed response using post service logic
		resp, err := s.postService.GetByID(ctx, p.ID)
		if err == nil {
			responses = append(responses, *resp)
		}
	}

	var nextCursor string
	if len(posts) > 0 && hasMore {
		last := posts[len(posts)-1]
		nextCursor = pagination.EncodeCursor(last.CreatedAt, last.ID.String())
	}

	return &BookmarkListResponse{
		Posts:   responses,
		Cursor:  nextCursor,
		HasMore: hasMore,
	}, nil
}
