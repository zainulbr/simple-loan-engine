package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService is interface for token service
type TokenService interface {
	GenerateToken(userId, role string, timestamp time.Time) string
	ValidateToken(token string) (*jwt.Token, error)
}
