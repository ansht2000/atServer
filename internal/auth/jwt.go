package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidOrExpiredToken = errors.New("token is invalid or expired")
var ErrRetrievingUserIDFromToken = errors.New("error error retrieving user id from token")
var ErrParsingUUIDFromString = errors.New("error parsing uuid id from string")

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "Chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject: userID.String(),
	})
	tokString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", errors.New("error signing token")
	}
	return tokString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, ErrInvalidOrExpiredToken
	}

	userString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, ErrRetrievingUserIDFromToken
	}

	userUUID, err := uuid.Parse(userString)
	if err != nil {
		return uuid.Nil, ErrParsingUUIDFromString
	}

	return userUUID, nil
}