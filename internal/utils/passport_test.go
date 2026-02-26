package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mytestpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	assert.NotEqual(t, password, hashedPassword, "Hashed password should not be the same as the original password")
}

func TestCheckPassword(t *testing.T) {
	password := "mytestpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	assert.True(t, CheckPasswordHash(password, hashedPassword), "CheckPasswordHash should return true for correct password")
	assert.False(t, CheckPasswordHash("wrongpassword", hashedPassword), "CheckPasswordHash should return false for incorrect password")
}
