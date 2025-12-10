package application

import "google.golang.org/grpc/codes"

type ErrorCode string

const (
	CodeUserServiceCallFailed ErrorCode = "user_service_call_failed"
	CodeNotFound              ErrorCode = "not_found"
	CodeInvalidArgument       ErrorCode = "invalid_argument"
	CodeNotAuthorizedError    ErrorCode = "not_authorized"
)

type ApplicationError interface {
	error
	ConventionalCode() ErrorCode
	GrpcCode() codes.Code
}

type applicationError struct {
	errorCode ErrorCode
	msg       string
	grpcCode  codes.Code
}

func (e *applicationError) Error() string {
	return e.msg
}

func (e *applicationError) ConventionalCode() ErrorCode {
	return e.errorCode
}

func (e *applicationError) GrpcCode() codes.Code {
	return e.grpcCode
}

func NewUserServiceCallFailedError(msg string) ApplicationError {
	return &applicationError{
		errorCode: CodeUserServiceCallFailed,
		msg:       msg,
		grpcCode:  codes.Unavailable,
	}
}

func NewNotFoundError(msg string) ApplicationError {
	return &applicationError{
		errorCode: CodeNotFound,
		msg:       msg,
		grpcCode:  codes.NotFound,
	}
}

func NewInvalidArgumentError(msg string) ApplicationError {
	return &applicationError{
		errorCode: CodeInvalidArgument,
		msg:       msg,
		grpcCode:  codes.InvalidArgument,
	}
}

func NewNotAuthorizedError(msg string) ApplicationError {
	return &applicationError{
		errorCode: CodeNotAuthorizedError,
		msg:       msg,
		grpcCode:  codes.Internal,
	}
}
