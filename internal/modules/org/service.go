package org

import (
	"context"
	"strings"
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
	members, err := s.repo.GetMembersWithUsers(ctx, orgID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	res := make([]OrgMemberResponse, 0, len(members))
	for _, m := range members {
		displayName := ""
		if m.DisplayName != nil {
			displayName = *m.DisplayName
		}
		avatarURL := ""
		if m.AvatarURL != nil {
			avatarURL = *m.AvatarURL
		}
		res = append(res, OrgMemberResponse{
			UserID:      m.UserID.String(),
			Username:    m.Username,
			DisplayName: displayName,
			AvatarURL:   avatarURL,
			Role:        m.Role,
			JoinedAt:    m.JoinedAt,
		})
	}
	return res, nil
}

func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID) (*OrgResponse, error) {
	org, err := s.repo.GetByID(ctx, orgID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if org == nil {
		return nil, apperrors.NotFound("Organization")
	}
	return s.toResponse(org), nil
}

func (s *Service) GetBySlug(ctx context.Context, slugValue string) (*OrgResponse, error) {
	org, err := s.repo.GetBySlug(ctx, slugValue)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if org == nil {
		return nil, apperrors.NotFound("Organization")
	}
	return s.toResponse(org), nil
}

func (s *Service) List(ctx context.Context, page, limit int) (*OrgListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	orgs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	res := make([]OrgResponse, 0, len(orgs))
	for _, o := range orgs {
		org := o
		res = append(res, *s.toResponse(&org))
	}

	hasMore := int64(offset+len(orgs)) < total
	return &OrgListResponse{
		Organizations: res,
		Page:          page,
		Limit:         limit,
		Total:         total,
		HasMore:       hasMore,
	}, nil
}

func (s *Service) UpdateOrg(ctx context.Context, actorID, orgID uuid.UUID, req *UpdateOrgRequest) (*OrgResponse, error) {
	role, err := s.repo.GetUserRole(ctx, orgID, actorID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if role != "owner" && role != "admin" {
		return nil, apperrors.Forbidden("You don't have permission to update this organization")
	}

	org, err := s.repo.GetByID(ctx, orgID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if org == nil {
		return nil, apperrors.NotFound("Organization")
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name != "" {
			org.Name = name
			newSlug := slug.Make(name)
			if newSlug != "" && newSlug != org.Slug {
				existing, _ := s.repo.GetBySlug(ctx, newSlug)
				if existing != nil && existing.ID != org.ID {
					newSlug = newSlug + "-" + uuid.New().String()[:8]
				}
				org.Slug = newSlug
			}
		}
	}

	if req.Bio != nil {
		org.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		org.AvatarURL = *req.AvatarURL
	}
	if req.WebsiteURL != nil {
		org.WebsiteURL = *req.WebsiteURL
	}
	org.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, org); err != nil {
		return nil, apperrors.Internal(err)
	}
	return s.toResponse(org), nil
}

func (s *Service) DeleteOrg(ctx context.Context, actorID, orgID uuid.UUID) error {
	role, err := s.repo.GetUserRole(ctx, orgID, actorID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if role != "owner" {
		return apperrors.Forbidden("Only the owner can delete this organization")
	}
	if err := s.repo.Delete(ctx, orgID); err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (s *Service) RemoveMember(ctx context.Context, actorID, orgID, targetUserID uuid.UUID) error {
	role, err := s.repo.GetUserRole(ctx, orgID, actorID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if role != "owner" && role != "admin" {
		return apperrors.Forbidden("You don't have permission to remove members")
	}

	memberRole, err := s.repo.GetUserRole(ctx, orgID, targetUserID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if memberRole == "owner" {
		return apperrors.Forbidden("Cannot remove the organization owner")
	}
	if memberRole == "" {
		return apperrors.NotFound("Member")
	}

	if err := s.repo.RemoveMember(ctx, orgID, targetUserID); err != nil {
		return apperrors.Internal(err)
	}
	return nil
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
