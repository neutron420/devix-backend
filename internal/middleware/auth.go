package middleware

import (
	"strings"

	apperrors "devix-backend/internal/errors"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextKeyUserID is the gin context key for the authenticated user's ID.
	ContextKeyUserID   = "user_id"
	// ContextKeyUsername is the gin context key for the authenticated user's username.
	ContextKeyUsername = "username"
	// ContextKeyUserRole is the gin context key for the authenticated user's role.
	ContextKeyUserRole = "user_role"
)

// Auth returns a middleware that validates JWT access tokens.
func Auth(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Abort(c, apperrors.Unauthorized("Authorization header is required"))
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Abort(c, apperrors.Unauthorized("Invalid authorization header format"))
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			response.Abort(c, apperrors.Unauthorized("Token is required"))
			return
		}

		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.Abort(c, apperrors.Unauthorized("Invalid or expired token"))
			return
		}

		// Set user context values for downstream handlers
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyUserRole, claims.Role)

		c.Next()
	}
}

// OptionalAuth is like Auth but allows unauthenticated requests to pass through.
func OptionalAuth(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(strings.TrimSpace(parts[1]))
		if err != nil {
			c.Next()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyUserRole, claims.Role)

		c.Next()
	}
}

// RequireRole returns a middleware that checks the user has one of the specified roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyUserRole)
		if !exists {
			response.Abort(c, apperrors.Unauthorized("Authentication required"))
			return
		}

		userRole := role.(string)
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		response.Abort(c, apperrors.Forbidden("Insufficient permissions"))
	}
}

// GetUserID extracts the authenticated user's ID from the gin context.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextKeyUserID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// GetUsername extracts the authenticated user's username from the gin context.
func GetUsername(c *gin.Context) string {
	val, _ := c.Get(ContextKeyUsername)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
