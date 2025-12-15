package domain_test

import (
	"testing"

	domain "github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"google.golang.org/grpc/codes"
)

type dErr interface {
	ConventionalCode() domain.ErrorCode
	Error() string
	GrpcCode() codes.Code
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		ctor     func(string) dErr
		msg      string
		wantCode domain.ErrorCode
		wantMsg  string
		wantGrpc codes.Code
	}{
		{
			name:     "NewFailedInvariantError",
			ctor:     func(m string) dErr { return domain.NewFailedInvariantError(m) },
			msg:      "test",
			wantCode: domain.CodeFailedInvariant,
			wantMsg:  "test",
			wantGrpc: codes.InvalidArgument,
		},
		{
			name:     "NewInvalidArgumentError",
			ctor:     func(m string) dErr { return domain.NewInvalidArgumentError(m) },
			msg:      "test",
			wantCode: domain.CodeInvalidArgument,
			wantMsg:  "test",
			wantGrpc: codes.InvalidArgument,
		},
		{
			name:     "NewInternalError",
			ctor:     func(m string) dErr { return domain.NewInternalError(m) },
			msg:      "test",
			wantCode: domain.CodeInternalError,
			wantMsg:  "test",
			wantGrpc: codes.Internal,
		},
		{
			name:     "NewAuthInvalidTokenError",
			ctor:     func(m string) dErr { return domain.NewAuthInvalidTokenError(m) },
			msg:      "test",
			wantCode: domain.CodeAuthInvalidToken,
			wantMsg:  "test",
			wantGrpc: codes.Unauthenticated,
		},
		{
			name:     "NewAuthExpiredTokenError",
			ctor:     func(m string) dErr { return domain.NewAuthExpiredTokenError(m) },
			msg:      "test",
			wantCode: domain.CodeAuthExpiredToken,
			wantMsg:  "test",
			wantGrpc: codes.Unauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ctor(tc.msg)

			if got := err.ConventionalCode(); got != tc.wantCode {
				t.Fatalf("ConventionalCode(): got %v, want %v", got, tc.wantCode)
			}
			if got := err.Error(); got != tc.wantMsg {
				t.Fatalf("Error(): got %q, want %q", got, tc.wantMsg)
			}
			if got := err.GrpcCode(); got != tc.wantGrpc {
				t.Fatalf("GrpcCode(): got %v, want %v", got, tc.wantGrpc)
			}
		})
	}
}
