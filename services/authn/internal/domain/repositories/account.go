package repositories

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type AccountRepository interface {
	NextIdentity() (*vo.AccountID, error)
	Save(ctx context.Context, account *aggregates.Account) error
	FindByID(ctx context.Context, account *vo.AccountID) (*aggregates.Account, error)
	FindByEmail(ctx context.Context, email *vo.Email) (*aggregates.Account, error)
	FindByLogin(ctx context.Context, login string) (*aggregates.Account, error)
}
