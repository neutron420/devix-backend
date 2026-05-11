package auth

import (
	"fmt"
	"unicode"

	"devix-backend/internal/config"
)

func ValidatePassword(policy config.PasswordPolicyConfig, password string) error {
	if policy.MinLength > 0 && len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}

	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSymbol bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsNumber(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	if policy.RequireUpper && !hasUpper {
		return fmt.Errorf("password must include at least one uppercase letter")
	}
	if policy.RequireLower && !hasLower {
		return fmt.Errorf("password must include at least one lowercase letter")
	}
	if policy.RequireNumber && !hasNumber {
		return fmt.Errorf("password must include at least one number")
	}
	if policy.RequireSymbol && !hasSymbol {
		return fmt.Errorf("password must include at least one symbol")
	}

	return nil
}
