package grpc

import (
	"context"
	"errors"
	"log"

	common "github.com/bcchicr/cinemaroom/services/authn/generated/cinemaroom/common/v1"
	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const statusCodePrefix = "cinemaroom_authn_"

type grpcMappableError[Code ~string] interface {
	error
	ConventionalCode() Code
	GrpcCode() codes.Code
}

func NewErrorHandlerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		if mapped, ok := mapError[domain.DomainError, domain.ErrorCode](
			err,
			info,
			"domain_",
		); ok {
			return nil, mapped
		}

		if mapped, ok := mapError[application.ApplicationError, application.ErrorCode](
			err,
			info,
			"application_",
		); ok {
			return nil, mapped
		}

		if mapped, ok := mapError[infrastructure.InfrastructureError, infrastructure.ErrorCode](
			err,
			info,
			"infrastructure_",
		); ok {
			return nil, mapped
		}

		log.Printf(
			"error: %s (unknown): %s", info.FullMethod, err.Error(),
		)

		return nil, buildStatusError(
			info.FullMethod,
			codes.Internal,
			statusCodePrefix+"unknown",
			"internal error",
		)
	}
}

func mapError[T grpcMappableError[Code], Code ~string](
	err error,
	info *grpc.UnaryServerInfo,
	prefix string,
) (mapped error, ok bool) {
	var target T
	if !errors.As(err, &target) {
		return nil, false
	}

	code := target.ConventionalCode()
	grpcCode := target.GrpcCode()
	customCode := statusCodePrefix + prefix + string(code)

	return buildStatusError(
		info.FullMethod,
		grpcCode,
		customCode,
		target.Error(),
	), true
}

func buildStatusError(
	method string,
	grpcCode codes.Code,
	customCode string,
	msg string,
) error {
	log.Printf(
		"error: %s (%s): %s",
		method,
		customCode,
		msg,
	)

	st := status.New(grpcCode, msg)
	details := &common.CustomErrorDetails{
		Code:    customCode,
		Message: msg,
	}

	stWithDetails, err := st.WithDetails(details)
	if err != nil {
		log.Printf(
			"failed to attach details for %s (%s): %v",
			method,
			customCode,
			err,
		)
		return st.Err()
	}

	return stWithDetails.Err()
}
