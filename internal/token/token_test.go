package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret"
	issuer := "test-issuer"
	service := NewService(secret, issuer)

	clientID := "test-client"
	duration := time.Minute

	// Test Generation
	token, err := service.GenerateAccessToken(clientID, duration)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Test Validation
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims["sub"] != clientID {
		t.Errorf("Expected sub %s, got %v", clientID, claims["sub"])
	}
	if claims["iss"] != issuer {
		t.Errorf("Expected iss %s, got %v", issuer, claims["iss"])
	}
}

func TestInvalidToken(t *testing.T) {
	service := NewService("secret", "issuer")
	_, err := service.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestExpiredToken(t *testing.T) {
	service := NewService("secret", "issuer")
	// Generate a token that expired 1 second ago
	claims := jwt.MapClaims{
		"sub": "client",
		"exp": time.Now().Add(-time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))

	_, err := service.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}
