package gateway

import (
	"context"
	"fmt"

	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGrpcClientErrorsInterceptor() grpc.UnaryClientInterceptor {

	return func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}

		st, ok := status.FromError(err)
		if !ok {
			return infrastructure.NewGrpcRequestFailedError(
				fmt.Sprintf("%s: non-grpc error: %v", method, err),
			)
		}

		switch st.Code() {
		case codes.Unavailable, codes.DeadlineExceeded, codes.Canceled:
			return infrastructure.NewExternalServiceDownError(
				fmt.Sprintf("%s: service unavailable (%s): %s", method, st.Code(), st.Message()),
			)

		default:
			return err
		}
	}
}
