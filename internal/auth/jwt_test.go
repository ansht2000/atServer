package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTValidation(t *testing.T) {
	userIDs := []uuid.UUID{uuid.New(), uuid.New()}

	var tokenSecret = "secret"

	cases := []struct{
		userID uuid.UUID
		tokenSecret string
		expiresIn time.Duration
		expectedID uuid.UUID
		expectedError error
	}{
		{
			userID: userIDs[0],
			tokenSecret: "secret",
			expiresIn: time.Hour,
			expectedID: userIDs[0],
			expectedError: nil,
		},
		{
			userID: userIDs[1],
			tokenSecret: "secret",
			expiresIn: time.Hour,
			expectedID: userIDs[1],
			expectedError: nil,
		},
		{
			userID: userIDs[0],
			tokenSecret: "secret",
			expiresIn: time.Nanosecond,
			expectedID: uuid.UUID{},
			expectedError: ErrInvalidOrExpiredToken,
		},
		{
			userID: userIDs[0],
			tokenSecret: "not_secret",
			expiresIn: time.Hour,
			expectedID: uuid.UUID{},
			expectedError: ErrInvalidOrExpiredToken,
		},

	}

	for _, c := range cases {
		tokString, err := MakeJWT(c.userID, c.tokenSecret, c.expiresIn)
		if err != nil {
			t.Errorf("Test failed due to signing of a token failing: %v", err)
			t.Fail()
		}
		// Sleep for just longer than the expiry time for one of the cases
		time.Sleep(2 * time.Nanosecond)

		userID, err := ValidateJWT(tokString, tokenSecret)
		if err != c.expectedError || userID != c.expectedID {
			t.Errorf("Failed to validate token: %v", err)
			t.Fail()
		}
	}
}