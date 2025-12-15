package handlers

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/application/clients"
	"github.com/bcchicr/cinemaroom/services/authn/internal/application/handlers/resources"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/repositories"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/password"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/token"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type RegisterCommand struct {
	Login    string
	Email    *string
	Password string
}

type RegistrationHandler interface {
	Handle(ctx context.Context, command RegisterCommand) (*vo.AccessToken, *resources.RefreshTokenResource, error)
}

type registrationHandler struct {
	passwordService        password.Service
	tokenService           token.Service
	accountRepository      repositories.AccountRepository
	refreshTokenRepository repositories.RefreshTokenRepository
	userServiceClient      clients.UserServiceClient
}

func NewRegistrationHandler(
	passwordService password.Service,
	tokenService token.Service,
	accountRepository repositories.AccountRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
	userServiceClient clients.UserServiceClient,
) RegistrationHandler {
	return &registrationHandler{
		passwordService:        passwordService,
		tokenService:           tokenService,
		accountRepository:      accountRepository,
		refreshTokenRepository: refreshTokenRepository,
		userServiceClient:      userServiceClient,
	}
}

func (handler *registrationHandler) Handle(ctx context.Context, command RegisterCommand) (*vo.AccessToken, *resources.RefreshTokenResource, error) {
	var email *vo.Email
	var err error
	if command.Email != nil {
		email, err = vo.NewEmail(*command.Email)
		if err != nil {
			return nil, nil, err
		}
	}

	a, err := handler.accountRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if a != nil {
		return nil, nil, application.NewNotAuthorizedError("email already taken")
	}

	a, err = handler.accountRepository.FindByLogin(ctx, command.Login)
	if err != nil {
		return nil, nil, err
	}
	if a != nil {
		return nil, nil, application.NewNotAuthorizedError("login already taken")
	}

	id, err := handler.accountRepository.NextIdentity()
	if err != nil {
		return nil, nil, err
	}

	password, err := vo.NewPassword(command.Password)
	if err != nil {
		return nil, nil, err
	}

	hash, err := handler.passwordService.Hash(password)
	if err != nil {
		return nil, nil, err
	}

	account, err := aggregates.NewAccount(
		id,
		command.Login,
		email,
		hash,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := handler.userServiceClient.Register(ctx, account.ID().Value(), account.Login()); err != nil {
		return nil, nil, err
	}

	err = handler.accountRepository.Save(ctx, account)
	if err != nil {
		return nil, nil, err
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
