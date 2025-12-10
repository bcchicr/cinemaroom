package infrastructure

import "google.golang.org/grpc/codes"

type ErrorCode string

const (
	CodeDBConnectionFailed  ErrorCode = "db_connection_failed"
	CodeDBQueryFailed       ErrorCode = "db_query_failed"
	CodeDBTransactionFailed ErrorCode = "db_transaction_failed"

	CodeGrpcConnectionFailed ErrorCode = "grpc_connection_failed"
	CodeGrpcRequestFailed    ErrorCode = "grpc_request_failed"

	CodeExternalServiceDown  ErrorCode = "external_service_down"
	CodeExternalServiceError ErrorCode = "external_service_error"

	CodeInternalError ErrorCode = "internal_error"
)

type InfrastructureError interface {
	error
	ConventionalCode() ErrorCode
	GrpcCode() codes.Code
}

type infrastructureError struct {
	errorCode ErrorCode
	msg       string
	grpcCode  codes.Code
}

func (e *infrastructureError) Error() string {
	return e.msg
}

func (e *infrastructureError) ConventionalCode() ErrorCode {
	return e.errorCode
}

func (e *infrastructureError) GrpcCode() codes.Code {
	return e.grpcCode
}

func NewDBConnectionFailedError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeDBConnectionFailed,
		msg:       msg,
		grpcCode:  codes.Unavailable,
	}
}

func NewDBQueryFailedError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeDBQueryFailed,
		msg:       msg,
		grpcCode:  codes.Internal,
	}
}

func NewDBTransactionFailedError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeDBTransactionFailed,
		msg:       msg,
		grpcCode:  codes.Aborted,
	}
}

func NewGrpcConnectionFailedError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeGrpcConnectionFailed,
		msg:       msg,
		grpcCode:  codes.Unavailable,
	}
}

func NewGrpcRequestFailedError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeGrpcRequestFailed,
		msg:       msg,
		grpcCode:  codes.Internal,
	}
}

func NewExternalServiceDownError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeExternalServiceDown,
		msg:       msg,
		grpcCode:  codes.Unavailable,
	}
}

func NewExternalServiceError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeExternalServiceError,
		msg:       msg,
		grpcCode:  codes.Internal,
	}
}

func NewInternalError(msg string) InfrastructureError {
	return &infrastructureError{
		errorCode: CodeInternalError,
		msg:       msg,
		grpcCode:  codes.Internal,
	}
}
