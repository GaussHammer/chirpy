package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("empty header")
	}
	authKey := strings.TrimPrefix(authHeader, "ApiKey ")
	return authKey, nil
}
