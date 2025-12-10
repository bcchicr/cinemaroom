package handlers

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/application/clients"
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
	Handle(context.Context, RegisterCommand) (*vo.AccessToken, *aggregates.RefreshToken, error)
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

func (handler *registrationHandler) Handle(ctx context.Context, command RegisterCommand) (*vo.AccessToken, *aggregates.RefreshToken, error) {
	id, err := handler.accountRepository.NextIdentity()
	if err != nil {
		return nil, nil, err
	}

	var email *vo.Email
	if command.Email != nil {
		var err error
		email, err = vo.NewEmail(*command.Email)
		if err != nil {
			return nil, nil, err
		}
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

	err = handler.refreshTokenRepository.Save(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
