package gateway

import (
	"context"
	"errors"
	"fmt"

	common "github.com/bcchicr/cinemaroom/services/authn/generated/cinemaroom/common/v1"
	user "github.com/bcchicr/cinemaroom/services/authn/generated/cinemaroom/user/v1"
	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type GrpcUserServiceClient struct {
	client user.UserServiceClient
}

func NewGrpcUserServiceClient(url string) (*GrpcUserServiceClient, error) {
	conn, err := grpc.NewClient(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			NewGrpcClientErrorsInterceptor(),
		),
	)
	if err != nil {
		return nil, infrastructure.NewGrpcConnectionFailedError(
			fmt.Sprintf("failed to connect to user service at %q: %v", url, err),
		)
	}

	client := user.NewUserServiceClient(conn)

	return &GrpcUserServiceClient{client: client}, nil
}

func (c *GrpcUserServiceClient) Register(ctx context.Context, userId uuid.UUID, username string) error {
	req := &user.RegisterRequest{
		Id: &common.Uuid{
			Value: userId.String(),
		},
		Username: username,
	}

	_, err := c.client.Register(ctx, req)
	if err == nil {
		return nil
	}

	var iErr infrastructure.InfrastructureError
	if errors.As(err, &iErr) {
		return iErr
	}

	st, ok := status.FromError(err)
	if !ok {
		return infrastructure.NewGrpcRequestFailedError(
			fmt.Sprintf("user.Register: unexpected non-grpc error: %v", err),
		)
	}

	for _, d := range st.Details() {
		switch info := d.(type) {
		case *common.CustomErrorDetails:
			return mapUserServiceCustomError(info)
		}
	}

	return infrastructure.NewExternalServiceError(
		fmt.Sprintf("user.Register error (%s): %s", st.Code(), st.Message()),
	)
}

func mapUserServiceCustomError(details *common.CustomErrorDetails) error {
	code := details.GetCode()
	msg := details.GetMessage()

	switch code {
	case "cinemaroom_users_domain_failed_invariant":
		return domain.NewFailedInvariantError(msg)

	case "cinemaroom_users_domain_invalid_argument":
		return domain.NewInvalidArgumentError(msg)

	case "cinemaroom_users_application_not_found":
		return application.NewNotFoundError(msg)

	case "cinemaroom_users_application_invalid_argument":
		return application.NewInvalidArgumentError(msg)

	default:
		return infrastructure.NewExternalServiceError(
			fmt.Sprintf("user service error %q: %s", code, msg),
		)
	}
}
