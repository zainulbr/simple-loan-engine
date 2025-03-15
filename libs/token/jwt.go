package token

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is a struct that will be encoded to a JWT Payload.
type Claims struct {
	UserId    string
	Role      string
	Timestamp time.Time
	jwt.RegisteredClaims
}

type jwtServices struct {
	secretKey string
	issure    string
}

// GenerateToken generates a new token
func (service *jwtServices) GenerateToken(userId, role string, timestamp time.Time) string {
	claims := &Claims{
		UserId:    userId,
		Role:      role,
		Timestamp: timestamp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1 hour
			Issuer:    service.issure,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

// ValidateToken validates a token
func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})

}

// NewService creates a new token service
func NewService(secretKey ...string) TokenService {
	_secretKey := getSecretKey()
	if len(secretKey) > 0 {
		_secretKey = secretKey[0]
	}
	return &jwtServices{
		secretKey: _secretKey,
	}
}

func getSecretKey() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "8D20809F3F8A1E8371B3D7DB989393A0E228B55EF3D845FF6B7F1A79C632FF80"
	}
	return secret
}
