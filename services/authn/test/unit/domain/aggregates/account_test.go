package aggregates

import (
	"errors"
	"testing"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/google/uuid"
)

func TestNewAccount(t *testing.T) {
	validAccountID, _ := vo.NewAccountIDFromString(uuid.NewString())
	validLogin := "test_login"
	validEmail, _ := vo.NewEmail("test@example.com")
	validPasswordHash, _ := vo.NewPasswordHash("hash")

	tests := []struct {
		name             string
		id               *vo.AccountID
		login            string
		email            *vo.Email
		passwordHash     *vo.PasswordHash
		wantErr          bool
		wantErrCode      domain.ErrorCode
		wantID           *vo.AccountID
		wantLogin        string
		wantEmail        *vo.Email
		wantPasswordHash *vo.PasswordHash
	}{
		{
			name:             "happy path",
			id:               validAccountID,
			login:            validLogin,
			email:            validEmail,
			passwordHash:     validPasswordHash,
			wantErr:          false,
			wantID:           validAccountID,
			wantLogin:        validLogin,
			wantEmail:        validEmail,
			wantPasswordHash: validPasswordHash,
		},
		{
			name:         "nil account id",
			id:           nil,
			login:        validLogin,
			email:        validEmail,
			passwordHash: validPasswordHash,
			wantErr:      true,
			wantErrCode:  domain.CodeInvalidArgument,
		},
		{
			name:         "empty login",
			id:           validAccountID,
			login:        "",
			email:        validEmail,
			passwordHash: validPasswordHash,
			wantErr:      true,
			wantErrCode:  domain.CodeInvalidArgument,
		},
		{
			name:             "empty email",
			id:               validAccountID,
			login:            validLogin,
			email:            nil,
			passwordHash:     validPasswordHash,
			wantErr:          false,
			wantID:           validAccountID,
			wantLogin:        validLogin,
			wantEmail:        nil,
			wantPasswordHash: validPasswordHash,
		},
		{
			name:         "empty password hash",
			id:           validAccountID,
			login:        validLogin,
			email:        validEmail,
			passwordHash: nil,
			wantErr:      true,
			wantErrCode:  domain.CodeInvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, err := aggregates.NewAccount(
				tc.id,
				tc.login,
				tc.email,
				tc.passwordHash,
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

			if got := a.Login(); got != tc.wantLogin {
				t.Fatalf("Login(): got %v, want %v", got, tc.wantLogin)
			}

			assertOptionalEqual(
				t,
				a.Email(),
				tc.wantEmail,
				"Email()",
				func(a, b *vo.Email) bool {
					return a.Equals(b)
				},
			)

			if got := a.PasswordHash(); !got.Equals(tc.wantPasswordHash) {
				t.Fatalf("PasswordHash(): got %v, want %v", got, tc.wantPasswordHash)
			}
		})
	}
}

func TestChangeLogin(t *testing.T) {
	validAccountID, _ := vo.NewAccountIDFromString(uuid.NewString())
	validEmail, _ := vo.NewEmail("test@example.com")
	validPasswordHash, _ := vo.NewPasswordHash("hash")

	tests := []struct {
		name        string
		oldLogin    string
		newLogin    string
		wantErr     bool
		wantErrCode domain.ErrorCode
		wantLogin   string
	}{
		{
			name:      "happy path",
			oldLogin:  "old_login",
			newLogin:  "new_login",
			wantErr:   false,
			wantLogin: "new_login",
		},
		{
			name:        "empty login",
			oldLogin:    "old_login",
			newLogin:    "",
			wantErr:     true,
			wantErrCode: domain.CodeInvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, _ := aggregates.NewAccount(
				validAccountID,
				tc.oldLogin,
				validEmail,
				validPasswordHash,
			)

			err := a.ChangeLogin(tc.newLogin)

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

			if got := a.Login(); got != tc.wantLogin {
				t.Fatalf("Login(): got %v, want %v", got, tc.wantLogin)
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	validAccountID, _ := vo.NewAccountIDFromString(uuid.NewString())
	validLogin := "login"
	validPasswordHash, _ := vo.NewPasswordHash("hash")

	someEmail, _ := vo.NewEmail("test1@test.com")
	otherEmail, _ := vo.NewEmail("test2@test.com")

	tests := []struct {
		name        string
		oldEmail    *vo.Email
		newEmail    *vo.Email
		wantErr     bool
		wantErrCode domain.ErrorCode
		wantEmail   *vo.Email
	}{
		{
			name:      "happy path",
			oldEmail:  someEmail,
			newEmail:  otherEmail,
			wantErr:   false,
			wantEmail: otherEmail,
		},
		{
			name:      "from nil",
			oldEmail:  nil,
			newEmail:  otherEmail,
			wantErr:   false,
			wantEmail: otherEmail,
		},
		{
			name:      "to nil",
			oldEmail:  someEmail,
			newEmail:  nil,
			wantErr:   false,
			wantEmail: nil,
		},
		{
			name:      "from nil to nil",
			oldEmail:  nil,
			newEmail:  nil,
			wantErr:   false,
			wantEmail: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, _ := aggregates.NewAccount(
				validAccountID,
				validLogin,
				tc.oldEmail,
				validPasswordHash,
			)

			err := a.ChangeEmail(tc.newEmail)

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

			assertOptionalEqual(
				t,
				a.Email(),
				tc.wantEmail,
				"Email()",
				func(a, b *vo.Email) bool {
					return a.Equals(b)
				},
			)
		})
	}
}

func TestChangePasswordHash(t *testing.T) {
	validAccountID, _ := vo.NewAccountIDFromString(uuid.NewString())
	validLogin := "login"
	validEmail, _ := vo.NewEmail("test@test.com")

	somePasswordHash, _ := vo.NewPasswordHash("hash1")
	otherPasswordHash, _ := vo.NewPasswordHash("hash2")

	tests := []struct {
		name             string
		oldPasswordHash  *vo.PasswordHash
		newPasswordHash  *vo.PasswordHash
		wantErr          bool
		wantErrCode      domain.ErrorCode
		wantPasswordHash *vo.PasswordHash
	}{
		{
			name:             "happy path",
			oldPasswordHash:  somePasswordHash,
			newPasswordHash:  otherPasswordHash,
			wantErr:          false,
			wantPasswordHash: otherPasswordHash,
		},
		{
			name:            "to nil",
			oldPasswordHash: somePasswordHash,
			newPasswordHash: nil,
			wantErr:         true,
			wantErrCode:     domain.CodeInvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, _ := aggregates.NewAccount(
				validAccountID,
				validLogin,
				validEmail,
				tc.oldPasswordHash,
			)

			err := a.ChangePasswordHash(tc.newPasswordHash)

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

			assertOptionalEqual(
				t,
				a.PasswordHash(),
				tc.wantPasswordHash,
				"PasswordHash()",
				func(a, b *vo.PasswordHash) bool {
					return a.Equals(b)
				},
			)
		})
	}
}
