package auth

import (
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if len(apiKey) == 0 {
		return "", ErrAuthorizationHeaderDoesNotExist
	}
	apiKey = strings.Split(apiKey, " ")[1]
	return apiKey, nil
}