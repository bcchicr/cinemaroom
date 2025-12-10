package repositories

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type AccountRepository interface {
	NextIdentity() (*vo.AccountID, error)
	Save(context.Context, *aggregates.Account) error
	FindByID(context.Context, *vo.AccountID) (*aggregates.Account, error)
	FindByEmail(context.Context, *vo.Email) (*aggregates.Account, error)
	FindByLogin(context.Context, string) (*aggregates.Account, error)
}
