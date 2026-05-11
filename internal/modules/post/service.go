package post

import (
	"context"
	"fmt"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/media"
	tagmod "devix-backend/internal/modules/tag"
	"devix-backend/internal/pkg/cache"
	"devix-backend/internal/pkg/pagination"
	"devix-backend/internal/pkg/sanitize"
	"devix-backend/internal/modules/search"
	"devix-backend/internal/queue"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo          *Repository
	mediaService  *media.Service
	tagService    *tagmod.Service
	cache         *cache.Cache
	queue         *queue.Queue
	followService interface {
		GetFollowingIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	}
	searchService *search.Service
	log           zerolog.Logger
}

func NewService(repo *Repository, mediaService *media.Service, tagService *tagmod.Service, searchService *search.Service, cache *cache.Cache, q *queue.Queue, log zerolog.Logger) *Service {
	s := &Service{
		repo:          repo,
		mediaService:  mediaService,
		tagService:    tagService,
		searchService: searchService,
		cache:         cache,
		queue:         q,
		log:           log.With().Str("module", "post").Logger(),
	}

	// Register job handlers
	if q != nil {
		q.Register("increment_view_count", s.handleIncrementViewCount)
		q.Register("sync_post_search", s.handleSyncPostSearch)
	}

	return s
}

func (s *Service) handleSyncPostSearch(ctx context.Context, payload interface{}) error {
	var postID uuid.UUID
	switch v := payload.(type) {
	case uuid.UUID:
		postID = v
	case string:
		var err error
		postID, err = uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid payload string for sync_post_search: %v", err)
		}
	default:
		return fmt.Errorf("invalid payload type for sync_post_search")
	}

	post, err := s.repo.GetByID(ctx, postID)
	if err != nil || post == nil {
		return err
	}

	tags, _ := s.repo.GetPostTags(ctx, postID)

	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = t.Name
	}

	return s.searchService.IndexPost(ctx, search.IndexedPost{
		ID:        post.ID.String(),
		Title:     post.Title,
		Content:   post.Content,
		Author:    post.AuthorID.String(),
		PostType:  string(post.PostType),
		Tags:      tagNames,
		CreatedAt: post.CreatedAt.UnixMilli(),
	})
}

func (s *Service) handleIncrementViewCount(ctx context.Context, payload interface{}) error {
	var id uuid.UUID
	switch v := payload.(type) {
	case uuid.UUID:
		id = v
	case string:
		var err error
		id, err = uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid payload string for increment_view_count: %v", err)
		}
	default:
		return fmt.Errorf("invalid payload type for increment_view_count")
	}
	return s.repo.IncrementViewCount(ctx, id)
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
			if err := s.repo.SetPostTags(ctx, post.ID, tagIDs); err != nil {
				s.log.Warn().Err(err).Msg("failed to set post tags")
			}
		}
	}
	if err := s.repo.IncrementUserPostCount(ctx, authorID); err != nil {
		s.log.Warn().Err(err).Msg("failed to increment user post count")
	}

	// Invalidate feed caches
	if err := s.cache.DeleteByPattern(ctx, "posts:feed:*"); err != nil {
		s.log.Warn().Err(err).Msg("failed to invalidate feed cache")
	}

	// Sync to Search (Elasticsearch) via Background Queue
	if s.queue != nil {
		s.queue.Enqueue(queue.Job{Type: "sync_post_search", Payload: post.ID})
	}

	return s.GetBySlug(ctx, slug)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*PostResponse, error) {
	cacheKey := fmt.Sprintf("posts:slug:%s", slug)
	var cached PostResponse
	if ok, _ := s.cache.Get(ctx, cacheKey, &cached); ok {
		if s.queue != nil {
			s.queue.Enqueue(queue.Job{Type: "increment_view_count", Payload: uuid.MustParse(cached.ID)})
		}
		return &cached, nil
	}

	post, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if post == nil {
		return nil, apperrors.NotFound("Post")
	}
	mediaList, _ := s.mediaService.GetPostMedia(ctx, post.ID)
	post.Media = mediaList
	if s.queue != nil {
		s.queue.Enqueue(queue.Job{Type: "increment_view_count", Payload: post.ID})
	}

	res := s.toResponse(post)
	_ = s.cache.Set(ctx, cacheKey, res, 30*time.Minute)
	return res, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*PostResponse, error) {
	cacheKey := fmt.Sprintf("posts:id:%s", id.String())
	var cached PostResponse
	if ok, _ := s.cache.Get(ctx, cacheKey, &cached); ok {
		return &cached, nil
	}

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

	res := s.toResponse(post)
	_ = s.cache.Set(ctx, cacheKey, res, 30*time.Minute)
	return res, nil
}

func (s *Service) List(ctx context.Context, query FeedQuery) (*PostListResponse, error) {
	// Generate cache key based on query
	cacheKey := fmt.Sprintf("posts:feed:%s:%s:%s:%s:%s", query.Sort, query.Type, query.Tag, query.AuthorID, query.Cursor)
	var cached PostListResponse
	if ok, _ := s.cache.Get(ctx, cacheKey, &cached); ok {
		return &cached, nil
	}

	// If it's a search query, use Elasticsearch
	if query.Search != "" && s.searchService != nil {
		offset := 0
		if query.Cursor != "" {
			// Search uses simple offset, so we decode cursor to get offset if possible
			// or just use 0 for now as a simple implementation.
		}
		esIDs, err := s.searchService.SearchPosts(ctx, query.Search, query.Limit, offset)
		if err == nil && len(esIDs) > 0 {
			posts, _, err := s.repo.GetByIDs(ctx, esIDs)
			if err == nil {
				responses := make([]PostResponse, 0, len(posts))
				for _, p := range posts {
					tags, _ := s.repo.GetPostTags(ctx, p.ID)
					p.Tags = tags
					responses = append(responses, *s.toResponse(&p))
				}
				hasMore := len(responses) == query.Limit
				return &PostListResponse{Posts: responses, HasMore: hasMore}, nil
			}
		}
	}

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

	res := &PostListResponse{Posts: responses, Cursor: cursor, HasMore: hasMore}
	// Cache feeds for shorter time
	_ = s.cache.Set(ctx, cacheKey, res, 5*time.Minute)

	return res, nil
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

func (s *Service) ListExplore(ctx context.Context, userID uuid.UUID, query FeedQuery) (*PostListResponse, error) {
	cacheKey := fmt.Sprintf("posts:explore:%s:%s:%s:%s", userID.String(), query.Type, query.Tag, query.Cursor)
	var cached PostListResponse
	if ok, _ := s.cache.Get(ctx, cacheKey, &cached); ok {
		return &cached, nil
	}

	if userID != uuid.Nil && s.followService != nil {
		followingIDs, err := s.followService.GetFollowingIDs(ctx, userID)
		if err == nil && len(followingIDs) > 0 {
			excludeIDs := make([]string, len(followingIDs))
			for i, id := range followingIDs {
				excludeIDs[i] = id.String()
			}
			query.ExcludeAuthorIDs = excludeIDs
		}
	}

	query.Sort = "trending"
	res, err := s.List(ctx, query)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, res, 10*time.Minute)
	return res, nil
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

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("posts:id:%s", postID.String()))
	_ = s.cache.Delete(ctx, fmt.Sprintf("posts:slug:%s", post.Slug))
	_ = s.cache.DeleteByPattern(ctx, "posts:feed:*")

	// Sync to Search
	if s.queue != nil {
		s.queue.Enqueue(queue.Job{Type: "sync_post_search", Payload: postID})
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

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("posts:id:%s", postID.String()))
	_ = s.cache.Delete(ctx, fmt.Sprintf("posts:slug:%s", post.Slug))
	_ = s.cache.DeleteByPattern(ctx, "posts:feed:*")

	// Sync to Search (Delete)
	if s.searchService != nil {
		_ = s.searchService.DeletePost(ctx, postID.String())
	}

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

func (s *Service) ListDrafts(ctx context.Context, authorID uuid.UUID) ([]PostResponse, error) {
	posts, err := s.repo.ListDrafts(ctx, authorID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	responses := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		responses = append(responses, *s.toResponse(&p))
	}
	return responses, nil
}

func (s *Service) Autosave(ctx context.Context, postID, userID uuid.UUID, req *AutosaveRequest) error {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if post == nil {
		return apperrors.NotFound("Post")
	}
	if post.AuthorID != userID {
		return apperrors.Forbidden("You can only autosave your own drafts")
	}
	if post.Status != "draft" {
		return apperrors.BadRequest("Only drafts can be autosaved")
	}
	return s.repo.Autosave(ctx, postID, req.Title, req.Content)
}
