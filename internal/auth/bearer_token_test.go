package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	cases := []struct{
		header http.Header
		expectedBearerToken string
		expectedError error
	}{
		{
			header: http.Header{"Authorization": []string{"Bearer token"}},
			expectedBearerToken: "token",
			expectedError: nil,
		},
		{
			header: http.Header{},
			expectedBearerToken: "",
			expectedError: ErrAuthorizationHeaderDoesNotExist,
		},
		{
			header: http.Header{"Authorization": []string{"Bearer  token"}},
			expectedBearerToken: "",
			expectedError: nil,
		},
	}

	for _, c := range cases {
		bearerToken, err := GetBearerToken(c.header)
		if err != c.expectedError || bearerToken != c.expectedBearerToken {
			t.Errorf("Failed to get bearer token for token '%s', with error: %v", c.expectedBearerToken, err)
			t.Fail()
		}
	}
}