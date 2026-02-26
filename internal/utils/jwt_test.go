package utils

import (
	"testing"
	"time"

	"github.com/veetmoradiya3628/go-shop/internal/config"
)

func TestGenerateTokenPairSuccess(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:              "mysecretkey",
		ExpiresIn:           15 * time.Minute,
		RefreshTokenExpires: 7 * 24 * time.Hour,
	}
	userID := uint(123)
	email := "test@example.com"
	role := "user"

	accessToken, refreshToken, err := GenerateTokenPair(cfg, userID, email, role)
	if err != nil {
		t.Fatalf("GenerateTokenPair returned an error: %v", err)
	}
	if accessToken == "" {
		t.Fatal("Expected access token to be generated, got empty string")
	}
	if refreshToken == "" {
		t.Fatal("Expected refresh token to be generated, got empty string")
	}
}
func TestGenerateTokenPairInvalidSecret(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:              "",
		ExpiresIn:           15 * time.Minute,
		RefreshTokenExpires: 7 * 24 * time.Hour,
	}

	_, _, err := GenerateTokenPair(cfg, 123, "test@example.com", "user")
	if err != nil {
		t.Fatalf("GenerateTokenPair failed with empty secret: %v", err)
	}
}

func TestValidateTokenSuccess(t *testing.T) {
	secret := "test-secret-key"
	cfg := &config.JWTConfig{
		Secret:              secret,
		ExpiresIn:           15 * time.Minute,
		RefreshTokenExpires: 7 * 24 * time.Hour,
	}

	accessToken, _, err := GenerateTokenPair(cfg, 123, "test@example.com", "admin")
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}

	claims, err := ValidateToken(accessToken, secret)
	if err != nil {
		t.Fatalf("ValidateToken returned an error: %v", err)
	}
	if claims.UserID != 123 {
		t.Errorf("Expected UserID 123, got %d", claims.UserID)
	}
	if claims.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", claims.Email)
	}
	if claims.Role != "admin" {
		t.Errorf("Expected role admin, got %s", claims.Role)
	}
}

func TestValidateTokenInvalidToken(t *testing.T) {
	_, err := ValidateToken("invalid.token.here", "test-secret-key")
	if err == nil {
		t.Fatal("ValidateToken should return an error for invalid token")
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	secret := "test-secret-key"
	cfg := &config.JWTConfig{
		Secret:              secret,
		ExpiresIn:           15 * time.Minute,
		RefreshTokenExpires: 7 * 24 * time.Hour,
	}

	accessToken, _, err := GenerateTokenPair(cfg, 123, "test@example.com", "user")
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}

	_, err = ValidateToken(accessToken, "wrong-secret")
	if err == nil {
		t.Fatal("ValidateToken should return an error for wrong secret")
	}
}

func TestValidateTokenEmptyToken(t *testing.T) {
	_, err := ValidateToken("", "test-secret-key")
	if err == nil {
		t.Fatal("ValidateToken should return an error for empty token")
	}
}
