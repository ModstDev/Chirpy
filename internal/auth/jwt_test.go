package auth

import (
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
