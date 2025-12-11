package vo_test

import (
	"errors"
	"maps"
	"testing"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	vo "github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantErr     bool
		wantErrCode domain.ErrorCode
		wantValue   string
		wantString  string
		wantMap     map[string]any
	}{
		{
			name:       "valid string",
			value:      "password",
			wantErr:    false,
			wantValue:  "password",
			wantString: "password",
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
		{
			name:        "short string",
			value:       "short",
			wantErr:     true,
			wantErrCode: domain.CodeFailedInvariant,
		},
		{
			name:        "long string",
			value:       "long_string_long_string_long_string_long_string_long_string_long_string_long_string_long_string_long_string_long_string",
			wantErr:     true,
			wantErrCode: domain.CodeFailedInvariant,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo, err := vo.NewPassword(tc.value)

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

			if got := vo.String(); got != tc.wantString {
				t.Fatalf("String(): got %v, want %v", got, tc.wantString)
			}

			if got := vo.ToMap(); !maps.Equal(got, tc.wantMap) {
				t.Fatalf("ToMap(): got %v, want %v", got, tc.wantMap)
			}
		})
	}
}

func TestPasswordEquals(t *testing.T) {
	tests := []struct {
		name       string
		value1     string
		value2     string
		wantEquals bool
	}{
		{
			name:       "equal",
			value1:     "password",
			value2:     "password",
			wantEquals: true,
		},
		{
			name:       "not equal value",
			value1:     "password1",
			value2:     "password2",
			wantEquals: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vo1, err := vo.NewPassword(
				tc.value1,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			vo2, err := vo.NewPassword(
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
