package comment

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/websocket"
	"devix-backend/internal/pkg/sanitize"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const maxCommentDepth = 3

type Service struct {
	repo      *Repository
	wsService *websocket.Service
	log       zerolog.Logger
}

func NewService(repo *Repository, wsService *websocket.Service, log zerolog.Logger) *Service {
	return &Service{
		repo:      repo,
		wsService: wsService,
		log:       log.With().Str("module", "comment").Logger(),
	}
}

func (s *Service) Create(ctx context.Context, postID, authorID uuid.UUID, req *CreateCommentRequest) (*CommentResponse, error) {
	content := sanitize.HTML(req.Content)
	depth := 0
	var parentID *uuid.UUID
	if req.ParentID != nil {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, apperrors.BadRequest("Invalid parent comment ID")
		}
		parentID = &pid
		parentDepth, err := s.repo.GetParentDepth(ctx, pid)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		if parentDepth < 0 {
			return nil, apperrors.NotFound("Parent comment")
		}
		depth = parentDepth + 1
		if depth > maxCommentDepth {
			return nil, apperrors.BadRequest("Maximum reply depth reached")
		}
	}
	now := time.Now()
	comment := &models.Comment{ID: uuid.New(), PostID: postID, AuthorID: authorID, ParentID: parentID, Content: content, Depth: depth, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, apperrors.Internal(err)
	}
	_ = s.repo.IncrementPostCommentCount(ctx, postID)

	res := &CommentResponse{
		ID: comment.ID.String(), PostID: postID.String(), Content: content, Depth: depth,
		CreatedAt: now.Format(time.RFC3339), UpdatedAt: now.Format(time.RFC3339),
	}
	s.wsService.BroadcastEvent(ctx, "new_comment", res)

	return res, nil
}

func (s *Service) GetByPostID(ctx context.Context, postID uuid.UUID) ([]CommentResponse, error) {
	comments, err := s.repo.GetByPostID(ctx, postID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return s.buildTree(comments), nil
}

func (s *Service) Update(ctx context.Context, commentID, userID uuid.UUID, req *UpdateCommentRequest) error {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if comment == nil {
		return apperrors.NotFound("Comment")
	}
	if comment.AuthorID != userID {
		return apperrors.Forbidden("You can only edit your own comments")
	}
	content := sanitize.HTML(req.Content)
	return s.repo.Update(ctx, commentID, content)
}

func (s *Service) Delete(ctx context.Context, commentID, userID uuid.UUID, userRole string) error {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if comment == nil {
		return apperrors.NotFound("Comment")
	}
	if comment.AuthorID != userID && userRole != "admin" && userRole != "moderator" {
		return apperrors.Forbidden("You can only delete your own comments")
	}
	return s.repo.SoftDelete(ctx, commentID)
}

func (s *Service) buildTree(flat []models.Comment) []CommentResponse {
	responseMap := make(map[uuid.UUID]*CommentResponse)
	var roots []CommentResponse
	for _, c := range flat {
		cr := CommentResponse{
			ID: c.ID.String(), PostID: c.PostID.String(), Content: c.Content, Depth: c.Depth,
			VoteCount: c.VoteCount, IsDeleted: c.IsDeleted, Replies: []CommentResponse{},
			CreatedAt: c.CreatedAt.Format(time.RFC3339), UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
		}
		if c.ParentID != nil {
			pid := c.ParentID.String()
			cr.ParentID = &pid
		}
		if c.Author != nil {
			cr.Author = &AuthorResponse{ID: c.Author.ID.String(), Username: c.Author.Username, DisplayName: c.Author.DisplayName, AvatarURL: c.Author.AvatarURL}
		}
		responseMap[c.ID] = &cr
	}
	for _, c := range flat {
		cr := responseMap[c.ID]
		if c.ParentID != nil {
			if parent, ok := responseMap[*c.ParentID]; ok {
				parent.Replies = append(parent.Replies, *cr)
				continue
			}
		}
		roots = append(roots, *cr)
	}
	if roots == nil {
		roots = []CommentResponse{}
	}
	return roots
}
