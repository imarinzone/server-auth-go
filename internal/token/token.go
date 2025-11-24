package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service handles token operations
type Service struct {
	secretKey []byte
	issuer    string
}

// NewService creates a new token service
func NewService(secretKey string, issuer string) *Service {
	return &Service{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateAccessToken creates a new JWT access token for a client
func (s *Service) GenerateAccessToken(clientID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": clientID,
		"iss": s.issuer,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates the token string and returns the claims
func (s *Service) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
