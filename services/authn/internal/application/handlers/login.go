package handlers

import (
	"context"
	"fmt"

	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/repositories"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/token"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type LoginIdentifierType int

const (
	LoginIdentifierEmail = iota
	LoginIdentifierLogin
)

var identifierName = map[LoginIdentifierType]string{
	LoginIdentifierEmail: "email",
	LoginIdentifierLogin: "login",
}

type LoginCommand struct {
	Identifier     string
	IdentifierType LoginIdentifierType
	Password       string
}

type LoginHandler interface {
	Handle(context.Context, LoginCommand) (*vo.AccessToken, *aggregates.RefreshToken, error)
}

type loginHandler struct {
	tokenService           token.Service
	accountRepository      repositories.AccountRepository
	refreshTokenRepository repositories.RefreshTokenRepository
}

func NewLoginHandler(
	tokenService token.Service,
	accountRepository repositories.AccountRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
) LoginHandler {
	return &loginHandler{
		tokenService:           tokenService,
		accountRepository:      accountRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (handler *loginHandler) Handle(ctx context.Context, command LoginCommand) (*vo.AccessToken, *aggregates.RefreshToken, error) {
	var account *aggregates.Account
	var err error
	switch command.IdentifierType {
	case LoginIdentifierEmail:
		email, err := vo.NewEmail(command.Identifier)
		if err != nil {
			return nil, nil, err
		}

		account, err = handler.accountRepository.FindByEmail(ctx, email)
		if err != nil {
			return nil, nil, err
		}

	case LoginIdentifierLogin:
		account, err = handler.accountRepository.FindByLogin(ctx, command.Identifier)
		if err != nil {
			return nil, nil, err
		}

	default:
		return nil, nil, application.NewInvalidArgumentError(
			fmt.Sprintf("invalid identifier type: %v", identifierName[command.IdentifierType]),
		)
	}

	if account == nil {
		return nil, nil, application.NewNotFoundError(
			fmt.Sprintf(
				"not found account with identifier %q, type %q",
				command.Identifier,
				identifierName[command.IdentifierType],
			),
		)
	}

	accessToken, refreshToken, err := handler.tokenService.Generate(account)
	if err != nil {
		return nil, nil, err
	}

	err = handler.refreshTokenRepository.Save(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
