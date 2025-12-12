package vo_test

import (
	"errors"
	"maps"
	"testing"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	vo "github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

func TestNewPasswordHash(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantErr     bool
		wantErrCode domain.ErrorCode
		wantValue   string
		wantMap     map[string]any
	}{
		{
			name:      "valid string",
			value:     "hash",
			wantErr:   false,
			wantValue: "hash",
			wantMap: map[string]any{
				"value": "",
			},
		},
		{
			name:        "empty string",
			value:       "",
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo, err := vo.NewPasswordHash(tc.value)

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

			if got := vo.Value(); got != tc.wantValue {
				t.Fatalf("Value(): got %q, want %q", got, tc.wantValue)
			}

			if got := vo.String(); got != "" {
				t.Fatalf("String(): got %v, want empty string", got)
			}

			if got := vo.ToMap(); !maps.Equal(got, tc.wantMap) {
				t.Fatalf("ToMap(): got %v, want %v", got, tc.wantMap)
			}
		})
	}
}

func TestPasswordHashEquals(t *testing.T) {
	tests := []struct {
		name       string
		value1     string
		value2     string
		wantEquals bool
	}{
		{
			name:       "equal",
			value1:     "hash",
			value2:     "hash",
			wantEquals: true,
		},
		{
			name:       "not equal value",
			value1:     "hash1",
			value2:     "hash2",
			wantEquals: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo1, err := vo.NewPasswordHash(
				tc.value1,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			vo2, err := vo.NewPasswordHash(
				tc.value2,
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
