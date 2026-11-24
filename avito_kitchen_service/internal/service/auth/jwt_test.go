package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTService_GenerateAndParse(t *testing.T) {
	service := NewJWTService("test-secret")

	token, err := service.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := service.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("UserID = %d, want 42", claims.UserID)
	}
}

func TestJWTService_InvalidToken(t *testing.T) {
	service := NewJWTService("test-secret")

	_, err := service.ParseToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWTService_WrongSecret(t *testing.T) {
	service := NewJWTService("test-secret")
	otherService := NewJWTService("other-secret")

	token, err := service.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = otherService.ParseToken(token)
	if err == nil {
		t.Fatal("expected error for token signed with another secret")
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	service := NewJWTService("test-secret")

	claims := Claims{
		UserID: 42,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	_, err = service.ParseToken(tokenString)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}
