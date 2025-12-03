package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	cleanedToken := strings.TrimSpace(strings.ReplaceAll(auth, "Bearer", ""))
	return cleanedToken, nil
}
