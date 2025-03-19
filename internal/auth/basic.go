package auth

import (
	"encoding/base64"
	"errors"
	"strings"
)

func ParseBasicAuth(authHeader string) (username, password string, err error) {
	if authHeader == "" {
		return "", "", errors.New("authorization header is empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", errors.New("invalid authorization header format")
	}

	payload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("invalid base64 encoding")
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", errors.New("invalid basic auth format")
	}

	return pair[0], pair[1], nil
}
