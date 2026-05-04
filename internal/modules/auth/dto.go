package auth

// --- Request DTOs ---

// SignupRequest represents the registration request body.
type SignupRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30,username"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest represents the token refresh request body.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// --- Response DTOs ---

// AuthResponse is the standard auth response after login/signup.
type AuthResponse struct {
	User   UserResponse `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

// UserResponse is the user data returned in auth responses.
type UserResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
	Role        string  `json:"role"`
}

// TokenResponse holds the token pair.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}
