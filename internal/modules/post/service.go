package post

import (
	"context"
	"fmt"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/media"
	tagmod "devix-backend/internal/modules/tag"
	"devix-backend/internal/pkg/pagination"
	"devix-backend/internal/pkg/sanitize"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo          *Repository
	mediaService  *media.Service
	tagService    *tagmod.Service
	followService interface {
		GetFollowingIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	}
	log zerolog.Logger
}

func NewService(repo *Repository, mediaService *media.Service, tagService *tagmod.Service, log zerolog.Logger) *Service {
	return &Service{repo: repo, mediaService: mediaService, tagService: tagService, log: log.With().Str("module", "post").Logger()}
}

func (s *Service) SetFollowService(fs interface {
	GetFollowingIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}) {
	s.followService = fs
}

func (s *Service) Create(ctx context.Context, authorID uuid.UUID, req *CreatePostRequest) (*PostResponse, error) {
	title := sanitize.Text(req.Title)
	content := sanitize.HTML(req.Content)
	slug := sanitize.Slug(title)
	if slug == "" {
		slug = uuid.New().String()[:8]
	}
	baseSlug := slug
	for i := 1; ; i++ {
		exists, err := s.repo.SlugExists(ctx, slug)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		if !exists {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
		if i > 100 {
			slug = fmt.Sprintf("%s-%s", baseSlug, uuid.New().String()[:8])
			break
		}
	}
	status := models.PostStatusPublished
	if req.Status == "draft" {
		status = models.PostStatusDraft
	}
	now := time.Now()
	post := &models.Post{
		ID:            uuid.New(),
		AuthorID:      authorID,
		Title:         title,
		Slug:          slug,
		Content:       content,
		PostType:      models.PostType(req.PostType),
		Status:        status,
		ExternalLinks: req.ExternalLinks,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.repo.Create(ctx, post); err != nil {
		return nil, apperrors.Internal(err)
	}
	if len(req.Tags) > 0 {
		tagIDs, err := s.tagService.FindOrCreateByNames(ctx, req.Tags)
		if err == nil {
			_ = s.repo.SetPostTags(ctx, post.ID, tagIDs)
		}
	}
	_ = s.repo.IncrementUserPostCount(ctx, authorID)
	return s.GetBySlug(ctx, slug)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*PostResponse, error) {
	post, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if post == nil {
		return nil, apperrors.NotFound("Post")
	}
	mediaList, _ := s.mediaService.GetPostMedia(ctx, post.ID)
	post.Media = mediaList
	go func() { _ = s.repo.IncrementViewCount(context.Background(), post.ID) }()
	return s.toResponse(post), nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*PostResponse, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if post == nil {
		return nil, apperrors.NotFound("Post")
	}
	tags, _ := s.repo.GetPostTags(ctx, id)
	post.Tags = tags
	mediaList, _ := s.mediaService.GetPostMedia(ctx, id)
	post.Media = mediaList
	return s.toResponse(post), nil
}

func (s *Service) List(ctx context.Context, query FeedQuery) (*PostListResponse, error) {
	posts, hasMore, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	responses := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		tags, _ := s.repo.GetPostTags(ctx, p.ID)
		p.Tags = tags
		responses = append(responses, *s.toResponse(&p))
	}
	var cursor string
	if len(posts) > 0 && hasMore {
		last := posts[len(posts)-1]
		cursor = pagination.EncodeCursor(last.CreatedAt, last.ID.String())
	}
	return &PostListResponse{Posts: responses, Cursor: cursor, HasMore: hasMore}, nil
}

func (s *Service) ListFollowing(ctx context.Context, userID uuid.UUID, query FeedQuery) (*PostListResponse, error) {
	if s.followService == nil {
		return nil, apperrors.Internal(fmt.Errorf("follow service not initialized"))
	}
	followingIDs, err := s.followService.GetFollowingIDs(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if len(followingIDs) == 0 {
		return &PostListResponse{Posts: []PostResponse{}, HasMore: false}, nil
	}

	authorIDs := make([]string, len(followingIDs))
	for i, id := range followingIDs {
		authorIDs[i] = id.String()
	}
	query.AuthorIDs = authorIDs

	return s.List(ctx, query)
}

func (s *Service) Update(ctx context.Context, postID, userID uuid.UUID, req *UpdatePostRequest) (*PostResponse, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if post == nil {
		return nil, apperrors.NotFound("Post")
	}
	if post.AuthorID != userID {
		return nil, apperrors.Forbidden("You can only edit your own posts")
	}
	var title, content, externalLinks *string
	if req.Title != nil {
		t := sanitize.Text(*req.Title)
		title = &t
	}
	if req.Content != nil {
		c := sanitize.HTML(*req.Content)
		content = &c
	}
	if req.ExternalLinks != nil {
		externalLinks = req.ExternalLinks
	}
	if err := s.repo.Update(ctx, postID, title, content, externalLinks, req.Status); err != nil {
		return nil, apperrors.Internal(err)
	}
	if req.Tags != nil {
		tagIDs, err := s.tagService.FindOrCreateByNames(ctx, req.Tags)
		if err == nil {
			_ = s.repo.SetPostTags(ctx, postID, tagIDs)
		}
	}
	return s.GetByID(ctx, postID)
}

func (s *Service) Delete(ctx context.Context, postID, userID uuid.UUID, userRole string) error {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if post == nil {
		return apperrors.NotFound("Post")
	}
	if post.AuthorID != userID && userRole != "admin" && userRole != "moderator" {
		return apperrors.Forbidden("You can only delete your own posts")
	}
	if err := s.repo.SoftDelete(ctx, postID); err != nil {
		return apperrors.Internal(err)
	}
	_ = s.repo.DecrementUserPostCount(ctx, post.AuthorID)
	return nil
}

func (s *Service) toResponse(p *models.Post) *PostResponse {
	resp := &PostResponse{
		ID: p.ID.String(), Title: p.Title, Slug: p.Slug, Content: p.Content,
		PostType: string(p.PostType), Status: string(p.Status),
		ViewCount: p.ViewCount, VoteCount: p.VoteCount, CommentCount: p.CommentCount,
		IsPinned: p.IsPinned, ExternalLinks: p.ExternalLinks,
		Tags: make([]TagResponse, 0), Media: make([]MediaResponse, 0),
		CreatedAt: p.CreatedAt.Format(time.RFC3339), UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
	if p.Author != nil {
		resp.Author = &AuthorResponse{ID: p.Author.ID.String(), Username: p.Author.Username, DisplayName: p.Author.DisplayName, AvatarURL: p.Author.AvatarURL}
	}
	for _, t := range p.Tags {
		resp.Tags = append(resp.Tags, TagResponse{ID: t.ID.String(), Name: t.Name, Slug: t.Slug})
	}
	for _, m := range p.Media {
		mr := MediaResponse{ID: m.ID.String(), FileURL: m.FileURL, FileType: string(m.FileType), FileSize: m.FileSize, MimeType: m.MimeType}
		if m.OriginalName != nil {
			mr.OriginalName = *m.OriginalName
		}
		resp.Media = append(resp.Media, mr)
	}
	return resp
}
