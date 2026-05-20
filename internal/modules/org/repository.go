package org

import (
	"context"
	"errors"
	"time"

	"devix-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type MemberWithUser struct {
	UserID      uuid.UUID
	Role        string
	JoinedAt    time.Time
	Username    string
	DisplayName *string
	AvatarURL   *string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(org).Error; err != nil {
			return err
		}
		// Creator becomes owner in OrgMember table too
		member := &models.OrgMember{
			ID:     uuid.New(),
			OrgID:  org.ID,
			UserID: org.OwnerID,
			Role:   "owner",
		}
		return tx.Create(member).Error
	})
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	err := r.db.WithContext(ctx).First(&org, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &org, err
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	var org models.Organization
	err := r.db.WithContext(ctx).First(&org, "slug = ?", slug).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &org, err
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]models.Organization, error) {
	var orgs []models.Organization
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orgs).Error
	return orgs, err
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Organization{}).Count(&count).Error
	return count, err
}

func (r *Repository) Update(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *Repository) AddMember(ctx context.Context, member *models.OrgMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *Repository) GetMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrgMember, error) {
	var members []models.OrgMember
	err := r.db.WithContext(ctx).Where("org_id = ?", orgID).Find(&members).Error
	return members, err
}

func (r *Repository) GetMembersWithUsers(ctx context.Context, orgID uuid.UUID) ([]MemberWithUser, error) {
	var members []MemberWithUser
	err := r.db.WithContext(ctx).
		Table("org_members").
		Select("org_members.user_id, org_members.role, org_members.joined_at, users.username, users.display_name, users.avatar_url").
		Joins("JOIN users ON users.id = org_members.user_id").
		Where("org_members.org_id = ?", orgID).
		Order("org_members.joined_at ASC").
		Scan(&members).Error
	return members, err
}

func (r *Repository) GetUserRole(ctx context.Context, orgID, userID uuid.UUID) (string, error) {
	var member models.OrgMember
	err := r.db.WithContext(ctx).Where("org_id = ? AND user_id = ?", orgID, userID).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return member.Role, err
}

func (r *Repository) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("org_id = ? AND user_id = ?", orgID, userID).Delete(&models.OrgMember{}).Error
}

func (r *Repository) Delete(ctx context.Context, orgID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("org_id = ?", orgID).Delete(&models.OrgMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Organization{}, orgID).Error
	})
}
