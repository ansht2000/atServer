package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("could not hash password")
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("passwords do not match")
	}
	return nil
}