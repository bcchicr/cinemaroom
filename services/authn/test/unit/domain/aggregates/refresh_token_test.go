package aggregates

import (
	"errors"
	"testing"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/google/uuid"
)

func TestNewRefreshToken(t *testing.T) {
	validRefreshTokenID, _ := vo.NewRefreshTokenIdFromString(uuid.NewString())
	validAccountID, _ := vo.NewAccountIDFromString(uuid.NewString())
	validValue := "value"
	validExpiresAt := time.Now()

	tests := []struct {
		name          string
		id            *vo.RefreshTokenID
		accountID     *vo.AccountID
		value         string
		expiresAt     time.Time
		wantErr       bool
		wantErrCode   domain.ErrorCode
		wantID        *vo.RefreshTokenID
		wantAccountID *vo.AccountID
		wantValue     string
		wantExpiresAt time.Time
	}{
		{
			name:          "happy path",
			id:            validRefreshTokenID,
			accountID:     validAccountID,
			value:         validValue,
			expiresAt:     validExpiresAt,
			wantErr:       false,
			wantID:        validRefreshTokenID,
			wantAccountID: validAccountID,
			wantValue:     validValue,
			wantExpiresAt: validExpiresAt,
		},
		{
			name:        "nil id",
			id:          nil,
			accountID:   validAccountID,
			value:       validValue,
			expiresAt:   validExpiresAt,
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
		{
			name:        "nil account id",
			id:          validRefreshTokenID,
			accountID:   nil,
			value:       validValue,
			expiresAt:   validExpiresAt,
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
		{
			name:        "empty value",
			id:          validRefreshTokenID,
			accountID:   validAccountID,
			value:       "",
			expiresAt:   validExpiresAt,
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, err := aggregates.NewRefreshToken(
				tc.id,
				tc.accountID,
				tc.value,
				tc.expiresAt,
			)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				var dErr domain.DomainError
				if !errors.As(err, &dErr) {
					t.Fatalf("expected domain err, got %v: ", dErr)
				}

				if got := dErr.ConventionalCode(); got != tc.wantErrCode {
					t.Fatalf("expected code %v, got %v", got, tc.wantErrCode)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := a.ID(); !got.Equals(tc.wantID) {
				t.Fatalf("ID(): got %q, want %q", got, tc.wantID)
			}

			if got := a.AccountID(); !got.Equals(tc.wantAccountID) {
				t.Fatalf("AccountID(): got %q, want %q", got, tc.wantAccountID)
			}

			if got := a.Value(); got != tc.wantValue {
				t.Fatalf("Login(): got %v, want %v", got, tc.wantValue)
			}

			if got := a.ExpiresAt(); !got.Equal(tc.wantExpiresAt) {
				t.Fatalf("ExpiresAt(): got %v, want %v", got, tc.wantExpiresAt)
			}

			if got := a.Hash(); got == a.Value() {
				t.Fatalf("Hash(): hash is same as value")
			}
		})
	}
}
