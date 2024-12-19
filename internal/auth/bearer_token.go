package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrAuthorizationHeaderDoesNotExist = errors.New("authorization header does not exist")

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	// If bearer token is "" return an error
	if len(bearerToken) == 0 {
		return "", ErrAuthorizationHeaderDoesNotExist
	}
	/* 
	   Bearer token in the header should be of the form: "Bearer {token_string}"
	   so split the string by one space and take the second element of the resulting slice
	   TODO: figure out if you want to accept bearer tokens with arbitrary number of spaces
	   between "Bearer" and the token and adjust accordingly
	*/
	bearerToken = strings.Split(bearerToken, " ")[1]
	return bearerToken, nil
}