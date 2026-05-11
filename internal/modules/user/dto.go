package user

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,max=100"`
	Bio         *string `json:"bio" binding:"omitempty,max=1000"`
	Username    *string `json:"username" binding:"omitempty,min=3,max=30,username"`
	WebsiteURL  *string `json:"website_url" binding:"omitempty,url,max=255"`
	GitHubURL   *string `json:"github_url" binding:"omitempty,max=255"`
	TwitterURL  *string `json:"twitter_url" binding:"omitempty,max=255"`
	Location    *string `json:"location" binding:"omitempty,max=100"`
	Preferences *string `json:"preferences" binding:"omitempty"`
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type ProfileResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	WebsiteURL  string  `json:"website_url"`
	GitHubURL   string  `json:"github_url"`
	TwitterURL  string  `json:"twitter_url"`
	Location    string  `json:"location"`
	Preferences string  `json:"preferences"`
	Role        string  `json:"role"`
	IsVerified  bool    `json:"is_verified"`
	PostCount   int      `json:"post_count"`
	Reputation  int      `json:"reputation"`
	Level       int      `json:"level"`
	Badges      []string `json:"badges"`
	CreatedAt   string   `json:"created_at"`
}

type PublicProfileResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	PostCount   int      `json:"post_count"`
	Reputation  int      `json:"reputation"`
	Level       int      `json:"level"`
	Badges      []string `json:"badges"`
	CreatedAt   string   `json:"created_at"`
}
