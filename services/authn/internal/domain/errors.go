package domain

import "google.golang.org/grpc/codes"

type ErrorCode string

const (
	CodeFailedInvariant  ErrorCode = "failed_invariant"
	CodeInvalidArgument  ErrorCode = "invalid_argument"
	CodeInternalError    ErrorCode = "internal_error"
	CodeAuthInvalidToken ErrorCode = "auth_invalid_token"
	CodeAuthExpiredToken ErrorCode = "auth_expired_token"
)

type DomainError interface {
	error
	ConventionalCode() ErrorCode
	GrpcCode() codes.Code
}

type domainError struct {
	errorCode ErrorCode
	msg       string
	grpcCode  codes.Code
}

func (e *domainError) Error() string {
	return e.msg
}

func (e *domainError) ConventionalCode() ErrorCode {
	return e.errorCode
}

func (e *domainError) GrpcCode() codes.Code {
	return e.grpcCode
}

func NewFailedInvariantError(msg string) DomainError {
	return &domainError{
		errorCode: CodeFailedInvariant,
		msg:       msg,
		grpcCode:  codes.InvalidArgument,
	}
}

func NewInvalidArgumentError(msg string) DomainError {
	return &domainError{
		errorCode: CodeInvalidArgument,
		msg:       msg,
		grpcCode:  codes.InvalidArgument,
	}
}

func NewInternalError(msg string) DomainError {
	return &domainError{
		errorCode: CodeInternalError,
		msg:       msg,
		grpcCode:  codes.Internal,
	}
}

func NewAuthInvalidTokenError(msg string) DomainError {
	return &domainError{
		errorCode: CodeAuthInvalidToken,
		msg:       msg,
		grpcCode:  codes.Unauthenticated,
	}
}

func NewAuthExpiredTokenError(msg string) DomainError {
	return &domainError{
		errorCode: CodeAuthExpiredToken,
		msg:       msg,
		grpcCode:  codes.Unauthenticated,
	}
}
