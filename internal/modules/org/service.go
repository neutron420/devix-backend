package org

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/rs/zerolog"
)

type Service struct {
	repo *Repository
	log  zerolog.Logger
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With().Str("module", "org").Logger(),
	}
}

func (s *Service) CreateOrg(ctx context.Context, ownerID uuid.UUID, req *CreateOrgRequest) (*OrgResponse, error) {
	orgSlug := slug.Make(req.Name)
	
	// Check if slug exists
	existing, _ := s.repo.GetBySlug(ctx, orgSlug)
	if existing != nil {
		orgSlug = orgSlug + "-" + uuid.New().String()[:8]
	}

	org := &models.Organization{
		ID:        uuid.New(),
		Name:      req.Name,
		Slug:      orgSlug,
		Bio:       req.Bio,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, org); err != nil {
		return nil, apperrors.Internal(err)
	}

	return s.toResponse(org), nil
}

func (s *Service) AddMember(ctx context.Context, actorID, orgID uuid.UUID, req *AddMemberRequest) error {
	// Check permission (must be owner or admin)
	role, err := s.repo.GetUserRole(ctx, orgID, actorID)
	if err != nil || (role != "owner" && role != "admin") {
		return apperrors.Forbidden("You don't have permission to add members")
	}

	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return apperrors.BadRequest("Invalid user ID")
	}

	// Check if already a member
	existingRole, _ := s.repo.GetUserRole(ctx, orgID, targetUserID)
	if existingRole != "" {
		return apperrors.Conflict("User is already a member")
	}

	member := &models.OrgMember{
		ID:       uuid.New(),
		OrgID:    orgID,
		UserID:   targetUserID,
		Role:     req.Role,
		JoinedAt: time.Now(),
	}

	if err := s.repo.AddMember(ctx, member); err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (s *Service) GetMembers(ctx context.Context, orgID uuid.UUID) ([]OrgMemberResponse, error) {
	members, err := s.repo.GetMembers(ctx, orgID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	res := make([]OrgMemberResponse, 0, len(members))
	for _, m := range members {
		res = append(res, OrgMemberResponse{
			UserID:   m.UserID.String(),
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}
	return res, nil
}

func (s *Service) toResponse(org *models.Organization) *OrgResponse {
	return &OrgResponse{
		ID:         org.ID.String(),
		Name:       org.Name,
		Slug:       org.Slug,
		Bio:        org.Bio,
		AvatarURL:  org.AvatarURL,
		WebsiteURL: org.WebsiteURL,
		OwnerID:    org.OwnerID.String(),
		CreatedAt:  org.CreatedAt,
	}
}
