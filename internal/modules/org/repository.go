package org

import (
	"context"
	"errors"

	"devix-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
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
