package handlers

import (
	"context"
	"fmt"

	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/repositories"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/password"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/token"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type EditCommand struct {
	AccessTokenString string
	Login             *string
	Email             *string
	Password          *string
}

type EditHandler interface {
	Handle(context.Context, EditCommand) (*vo.AccountID, error)
}

type editHandler struct {
	passwordService        password.Service
	tokenService           token.Service
	accountRepository      repositories.AccountRepository
	refreshTokenRepository repositories.RefreshTokenRepository
}

func NewEditHandler(
	passwordService password.Service,
	tokenService token.Service,
	accountRepository repositories.AccountRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
) EditHandler {
	return &editHandler{
		passwordService:        passwordService,
		tokenService:           tokenService,
		accountRepository:      accountRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (handler *editHandler) Handle(ctx context.Context, command EditCommand) (*vo.AccountID, error) {
	claims, err := handler.tokenService.GetClaims(ctx, command.AccessTokenString)
	if err != nil {
		return nil, err
	}

	accountID, err := vo.NewAccountIDFromString(claims.AccountID)
	if err != nil {
		return nil, err
	}

	account, err := handler.accountRepository.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, application.NewNotFoundError(
			fmt.Sprintf("not found account with id: %s", accountID.String()),
		)
	}

	if command.Login != nil {
		account.ChangeLogin(*command.Login)
	}

	if command.Email != nil {
		email, err := vo.NewEmail(*command.Email)
		if err != nil {
			return nil, err
		}

		account.ChangeEmail(email)
	}

	if command.Password != nil {
		password, err := vo.NewPassword(*command.Password)
		if err != nil {
			return nil, err
		}

		hash, err := handler.passwordService.Hash(password)
		if err != nil {
			return nil, err
		}

		account.ChangePasswordHash(hash)
	}

	err = handler.accountRepository.Save(ctx, account)
	if err != nil {
		return nil, err
	}

	return account.ID(), nil
}
