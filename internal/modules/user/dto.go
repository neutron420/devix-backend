package user

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,max=100"`
	Bio         *string `json:"bio" binding:"omitempty,max=1000"`
	Username    *string `json:"username" binding:"omitempty,min=3,max=30,username"`
}

type ProfileResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	Role        string  `json:"role"`
	IsVerified  bool    `json:"is_verified"`
	PostCount   int     `json:"post_count"`
	Reputation  int     `json:"reputation"`
	CreatedAt   string  `json:"created_at"`
}

type PublicProfileResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	PostCount   int     `json:"post_count"`
	Reputation  int     `json:"reputation"`
	CreatedAt   string  `json:"created_at"`
}
