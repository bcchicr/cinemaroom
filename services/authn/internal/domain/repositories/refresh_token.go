package repositories

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type RefreshTokenRepository interface {
	NextIdentity() (*vo.RefreshTokenID, error)
	Save(ctx context.Context, token *aggregates.RefreshToken) error
	FindByID(ctx context.Context, id *vo.RefreshTokenID) (*aggregates.RefreshToken, error)
	FindByHash(ctx context.Context, hash string) (*aggregates.RefreshToken, error)
	Delete(ctx context.Context, token *aggregates.RefreshToken) error
}
