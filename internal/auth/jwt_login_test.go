package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT_Success(t *testing.T) {
	id := uuid.New()
	tokenSecret := "are you sure of that?"
	expiresIn := 24 * time.Hour

	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	validatedUUID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if validatedUUID != id {
		t.Errorf("expected id %v, got %v", id, validatedUUID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Romeo"
	expiresIn := time.Duration(0)

	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Errorf("expected error validating expiring token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Zika"
	expiresIn := time.Duration(0)

	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(tokenString, "Zuka")
	if err == nil {
		t.Errorf("expected error validating a wrong secret, got nil")
	}
}
