package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	if auth == "" {
		return "", fmt.Errorf("authentication header not found")
	}

	apiKey := strings.TrimSpace(strings.ReplaceAll(auth, "ApiKey", ""))
	log.Printf("apiKey: %v ", apiKey)

	return apiKey, nil
}
