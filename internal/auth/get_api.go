package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("Authorization header format is not 'ApiKey <token>'")
	}
	token := strings.TrimPrefix(authHeader, "ApiKey ")
	return token, nil
}
