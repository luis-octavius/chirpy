package auth

import (
	"testing"
)

func TestCheckPasswordHash_Success(t *testing.T) {
	password := "hakuna matata"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if !match {
		t.Errorf("expected match to be %v, got %v", true, match)
	}
}

func TestCheckPasswordHash_Fail(t *testing.T) {
	password := "playstation"

	hashedPassword, err := HashPassword("pumba la pumba")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPassword returned error: %v", err)
	}

	if match {
		t.Errorf("expected match to be %v, got %v", false, match)
	}
}
