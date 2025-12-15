package vo_test

import (
	"errors"
	"maps"
	"testing"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	vo "github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/google/uuid"
)

func TestNewRefreshTokenIDFromString(t *testing.T) {
	validUuid, _ := uuid.NewV7()

	tests := []struct {
		name        string
		rawID       string
		wantErr     bool
		wantErrCode domain.ErrorCode
		wantValue   uuid.UUID
		wantString  string
		wantMap     map[string]any
	}{
		{
			name:       "valid raw id",
			rawID:      validUuid.String(),
			wantErr:    false,
			wantValue:  validUuid,
			wantString: validUuid.String(),
			wantMap: map[string]any{
				"value": validUuid.String(),
			},
		},
		{
			name:        "empty raw id",
			rawID:       "",
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
		{
			name:        "invalid raw id",
			rawID:       "invalid uuid",
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo, err := vo.NewRefreshTokenIdFromString(tc.rawID)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				var dErr domain.DomainError
				if !errors.As(err, &dErr) {
					t.Fatalf("expected domain err, got %v: ", dErr)
				}

				if got := dErr.ConventionalCode(); got != tc.wantErrCode {
					t.Fatalf("expected code %v, got %v", tc.wantErrCode, got)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := vo.Value(); got != tc.wantValue {
				t.Fatalf("Value(): got %v, want %v", got, tc.wantValue)
			}

			if got := vo.String(); got != tc.wantString {
				t.Fatalf("String(): got %q, want %q", got, tc.wantString)
			}

			if got := vo.ToMap(); !maps.Equal(got, tc.wantMap) {
				t.Fatalf("ToMap(): got %v, want %v", got, tc.wantMap)
			}
		})
	}
}

func TestRefreshTokenIDEquals(t *testing.T) {
	someRawID := uuid.NewString()
	otherRawID := uuid.NewString()

	tests := []struct {
		name       string
		rawID1     string
		rawID2     string
		wantEquals bool
	}{
		{
			name:       "equal",
			rawID1:     someRawID,
			rawID2:     someRawID,
			wantEquals: true,
		},
		{
			name:       "not equal value",
			rawID1:     someRawID,
			rawID2:     otherRawID,
			wantEquals: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo1, err := vo.NewRefreshTokenIdFromString(
				tc.rawID1,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			vo2, err := vo.NewRefreshTokenIdFromString(
				tc.rawID2,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := vo1.Equals(vo2); got != tc.wantEquals {
				t.Fatalf("Equals(): got %v, want %v", got, tc.wantEquals)
			}
		})
	}
}
