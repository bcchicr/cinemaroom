package handlers

import (
	"context"
	"fmt"

	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/application/handlers/resources"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/repositories"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/password"
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
	Handle(ctx context.Context, command LoginCommand) (*vo.AccessToken, *resources.RefreshTokenResource, error)
}

type loginHandler struct {
	tokenService           token.Service
	passwordService        password.Service
	accountRepository      repositories.AccountRepository
	refreshTokenRepository repositories.RefreshTokenRepository
}

func NewLoginHandler(
	tokenService token.Service,
	passwordService password.Service,
	accountRepository repositories.AccountRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
) LoginHandler {
	return &loginHandler{
		tokenService:           tokenService,
		passwordService:        passwordService,
		accountRepository:      accountRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (handler *loginHandler) Handle(ctx context.Context, command LoginCommand) (*vo.AccessToken, *resources.RefreshTokenResource, error) {
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
		return nil, nil, application.NewNotAuthorizedError(
			fmt.Sprintf(
				"not found account with identifier %q, type %q",
				command.Identifier,
				identifierName[command.IdentifierType],
			),
		)
	}

	pass, err := vo.NewPassword(command.Password)
	if err != nil {
		return nil, nil, err
	}

	if ok := handler.passwordService.Verify(pass, account.PasswordHash()); !ok {
		return nil, nil, application.NewNotAuthorizedError("invalid password")
	}

	accessToken, refreshToken, err := handler.tokenService.Generate(account)
	if err != nil {
		return nil, nil, err
	}

	if refreshToken.Value() == nil {
		return nil, nil, domain.NewFailedInvariantError("refresh token must have not nil value")
	}

	err = handler.refreshTokenRepository.Save(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	return accessToken,
		&resources.RefreshTokenResource{
			Value:     *refreshToken.Value(),
			ExpiresAt: refreshToken.ExpiresAt(),
		},
		nil
}
