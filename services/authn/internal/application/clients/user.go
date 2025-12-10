package clients

import (
	"context"

	"github.com/google/uuid"
)

type UserServiceClient interface {
	Register(ctx context.Context, userId uuid.UUID, username string) error
}
