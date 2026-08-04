package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()

	validSecret := "super-secret"
	wrongSecret := "wrong-secret"

	validToken, _ := MakeJWT(userID, validSecret, time.Hour)
	expiredToken, _ := MakeJWT(userID, validSecret, -time.Minute)

	tests := []struct {
		name       string
		token      string
		secret     string
		wantErr    bool
		expectedID uuid.UUID
	}{
		{
			name:       "Valid token",
			token:      validToken,
			secret:     validSecret,
			wantErr:    false,
			expectedID: userID,
		},
		{
			name:       "Wrong secret",
			token:      validToken,
			secret:     wrongSecret,
			wantErr:    true,
			expectedID: uuid.Nil,
		},
		{
			name:       "Expired token",
			token:      expiredToken,
			secret:     validSecret,
			wantErr:    true,
			expectedID: uuid.Nil,
		},
		{
			name:       "Malformed token",
			token:      "this-is-not-a-jwt",
			secret:     validSecret,
			wantErr:    true,
			expectedID: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ValidateJWT(tt.token, tt.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && id != tt.expectedID {
				t.Errorf("ValidateJWT() = %v, expected %v", id, tt.expectedID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
