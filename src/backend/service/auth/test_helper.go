package auth

import (
	"local/infra/repo"
)

// NewTestAuthService creates a new auth service instance for testing
// This allows tests to create service instances with custom repository and jwtSecret
func NewTestAuthService(repo *repo.Repository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// GetJWTClaimsType returns the JWTClaims type for testing
// This allows tests to create JWTClaims instances
type TestJWTClaims = JWTClaims

