package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var ErrCouldNotMakeRefreshToken = errors.New("could make new refresh token")

func MakeRefreshToken() (string, error) {
	randNum := make([]byte, 32)
	_, err := rand.Read(randNum)
	if err != nil {
		return "", ErrCouldNotMakeRefreshToken
	}
	return hex.EncodeToString(randNum), nil
}