package auth

import (
	"net/http"
	"testing"
)

func TestGetAPIKey_Success(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Authorization", "f271c81ff7084ee5b99a5091b42d486e")
	expected := "f271c81ff7084ee5b99a5091b42d486e"

	apiKey, err := GetAPIKey(req.Header)
	if err != nil {
		t.Fatalf("GetAPIKey returned error: %v", err)
	}

	if apiKey != expected {
		t.Errorf("expectev %v, got %v", expected, apiKey)
	}
}
