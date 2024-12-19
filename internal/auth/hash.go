package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrCouldNotHashPassword = errors.New("could not hash password")
var ErrPasswordsDontMatch = errors.New("passwords do not match")

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrCouldNotHashPassword
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrPasswordsDontMatch
	}
	return nil
}