package repositories

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type RefreshTokenRepository interface {
	NextIdentity() (*vo.RefreshTokenID, error)
	Save(context.Context, *aggregates.RefreshToken) error
	FindByID(context.Context, *vo.RefreshTokenID) (*aggregates.RefreshToken, error)
	FindByValue(context.Context, string) (*aggregates.RefreshToken, error)
	Delete(context.Context, *aggregates.RefreshToken) error
	DeleteByValue(context.Context, string) error
}
