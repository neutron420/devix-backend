package org

import "time"

type CreateOrgRequest struct {
	Name string `json:"name" binding:"required,max=100"`
	Bio  string `json:"bio" binding:"max=1000"`
}

type UpdateOrgRequest struct {
	Name       *string `json:"name" binding:"omitempty,max=100"`
	Bio        *string `json:"bio" binding:"omitempty,max=1000"`
	AvatarURL  *string `json:"avatar_url" binding:"omitempty"`
	WebsiteURL *string `json:"website_url" binding:"omitempty,url"`
}

type OrgResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Bio        string    `json:"bio"`
	AvatarURL  string    `json:"avatar_url"`
	WebsiteURL string    `json:"website_url"`
	OwnerID    string    `json:"owner_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrgMemberResponse struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=admin member"`
}
