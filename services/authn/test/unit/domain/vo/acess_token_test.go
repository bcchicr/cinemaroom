package vo_test

import (
	"errors"
	"maps"
	"testing"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	vo "github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

func date(s string) time.Time {
	tm, _ := time.Parse(time.DateOnly, s)
	return tm
}

func TestNewAccessToken(t *testing.T) {

	tests := []struct {
		name        string
		jti         string
		value       string
		expiresAt   time.Time
		wantErr     bool
		wantErrCode domain.ErrorCode
		wantJti     string
		wantValue   string
		wantExpiry  time.Time
		wantString  string
		wantMap     map[string]any
	}{
		{
			name:       "valid token",
			jti:        "1",
			value:      "access_token",
			expiresAt:  date("2025-01-01"),
			wantErr:    false,
			wantJti:    "1",
			wantValue:  "access_token",
			wantExpiry: date("2025-01-01"),
			wantString: "access_token",
			wantMap: map[string]any{
				"jti":       "1",
				"value":     "access_token",
				"expiresAt": date("2025-01-01"),
			},
		},
		{
			name:        "empty jti",
			jti:         "",
			value:       "access_token",
			expiresAt:   date("2025-01-01"),
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
		{
			name:        "empty value",
			jti:         "1",
			value:       "",
			expiresAt:   date("2025-01-01"),
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo, err := vo.NewAccessToken(tc.jti, tc.value, tc.expiresAt)

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

			if got := vo.Jti(); got != tc.wantJti {
				t.Fatalf("Jti(): got %v, want %v", got, tc.wantJti)
			}

			if got := vo.Value(); got != tc.wantValue {
				t.Fatalf("Value(): got %q, want %q", got, tc.wantValue)
			}

			if got := vo.ExpiresAt(); !got.Equal(tc.wantExpiry) {
				t.Fatalf("ExpiresAt(): got %v, want %v", got, tc.wantExpiry)
			}

			if got := vo.String(); got != tc.wantString {
				t.Fatalf("String(): got %v, want %v", got, tc.wantString)
			}

			if got := vo.ToMap(); !maps.Equal(got, tc.wantMap) {
				t.Fatalf("ToMap(): got %v, want %v", got, tc.wantMap)
			}
		})
	}
}

func TestAccessTokenEquals(t *testing.T) {
	tests := []struct {
		name       string
		jti1       string
		jti2       string
		value1     string
		value2     string
		expiresAt1 time.Time
		expiresAt2 time.Time
		wantEquals bool
	}{
		{
			name:       "equal",
			jti1:       "1",
			jti2:       "1",
			value1:     "value",
			value2:     "value",
			expiresAt1: date("2025-01-01"),
			expiresAt2: date("2025-01-01"),
			wantEquals: true,
		},
		{
			name:       "not equal jti",
			jti1:       "1",
			jti2:       "2",
			value1:     "value",
			value2:     "value",
			expiresAt1: date("2025-01-01"),
			expiresAt2: date("2025-01-01"),
			wantEquals: false,
		},
		{
			name:       "not equal value",
			jti1:       "1",
			jti2:       "1",
			value1:     "value1",
			value2:     "value2",
			expiresAt1: date("2025-01-01"),
			expiresAt2: date("2025-01-01"),
			wantEquals: false,
		},
		{
			name:       "not equal expiry",
			jti1:       "1",
			jti2:       "1",
			value1:     "value",
			value2:     "value",
			expiresAt1: date("2025-01-01"),
			expiresAt2: date("2025-01-02"),
			wantEquals: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo1, err := vo.NewAccessToken(
				tc.jti1,
				tc.value1,
				tc.expiresAt1,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			vo2, err := vo.NewAccessToken(
				tc.jti2,
				tc.value2,
				tc.expiresAt2,
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
