package auth

import (
	"net/http"
	"testing"
)

func TestBearerToken_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	req.Header.Set("Authorization", "Bearer secret")
	expected := "secret"

	bearerToken, err := GetBearerToken(req.Header)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}

	if bearerToken != expected {
		t.Errorf("expected %v, got %v", expected, bearerToken)
	}
}

func TestBearerToken_Fail(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	req.Header.Set("Authorization", "Bearer token123")
	expected := "token456"

	bearerToken, err := GetBearerToken(req.Header)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}

	if bearerToken == expected {
		t.Errorf("expected %v, got %v", expected, bearerToken)
	}
}
