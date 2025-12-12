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
	validValueHash := "hash"
	validExpiresAt := time.Now()

	emptyValue := ""

	tests := []struct {
		name          string
		id            *vo.RefreshTokenID
		accountID     *vo.AccountID
		value         *string
		valueHash     string
		expiresAt     time.Time
		wantErr       bool
		wantErrCode   domain.ErrorCode
		wantID        *vo.RefreshTokenID
		wantAccountID *vo.AccountID
		wantValue     *string
		wantValueHash string
		wantExpiresAt time.Time
	}{
		{
			name:          "happy path",
			id:            validRefreshTokenID,
			accountID:     validAccountID,
			value:         &validValue,
			valueHash:     validValueHash,
			expiresAt:     validExpiresAt,
			wantErr:       false,
			wantID:        validRefreshTokenID,
			wantAccountID: validAccountID,
			wantValue:     &validValue,
			wantValueHash: validValueHash,
			wantExpiresAt: validExpiresAt,
		},
		{
			name:          "happy path nil value",
			id:            validRefreshTokenID,
			accountID:     validAccountID,
			value:         nil,
			valueHash:     validValueHash,
			expiresAt:     validExpiresAt,
			wantErr:       false,
			wantID:        validRefreshTokenID,
			wantAccountID: validAccountID,
			wantValue:     nil,
			wantValueHash: validValueHash,
			wantExpiresAt: validExpiresAt,
		},
		{
			name:        "nil id",
			id:          nil,
			accountID:   validAccountID,
			value:       &validValue,
			valueHash:   validValueHash,
			expiresAt:   validExpiresAt,
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
		{
			name:        "nil account id",
			id:          validRefreshTokenID,
			accountID:   nil,
			value:       &validValue,
			valueHash:   validValueHash,
			expiresAt:   validExpiresAt,
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
		{
			name:        "empty value",
			id:          validRefreshTokenID,
			accountID:   validAccountID,
			value:       &emptyValue,
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
				tc.valueHash,
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
				t.Fatalf("Value(): got %v, want %v", got, tc.wantValue)
			}

			if got := a.ValueHash(); got != tc.wantValueHash {
				t.Fatalf("ValueHash(): got %v, want %v", got, tc.wantValueHash)
			}

			if got := a.ExpiresAt(); !got.Equal(tc.wantExpiresAt) {
				t.Fatalf("ExpiresAt(): got %v, want %v", got, tc.wantExpiresAt)
			}
		})
	}
}
