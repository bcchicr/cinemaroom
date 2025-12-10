package handlers

import (
	"context"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/repositories"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/token"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type RefreshCommand struct {
	RefreshTokenString string
}

type RefreshHandler interface {
	Handle(context.Context, RefreshCommand) (*vo.AccessToken, *aggregates.RefreshToken, error)
}

type refreshHandler struct {
	tokenService           token.Service
	accountRepository      repositories.AccountRepository
	refreshTokenRepository repositories.RefreshTokenRepository
}

func NewRefreshHandler(
	tokenService token.Service,
	accountRepository repositories.AccountRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
) RefreshHandler {
	return &refreshHandler{
		tokenService:           tokenService,
		accountRepository:      accountRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (handler *refreshHandler) Handle(ctx context.Context, command RefreshCommand) (*vo.AccessToken, *aggregates.RefreshToken, error) {
	refreshToken, err := handler.refreshTokenRepository.FindByValue(ctx, command.RefreshTokenString)
	if err != nil {
		return nil, nil, err
	}

	if refreshToken == nil {
		return nil, nil, domain.NewAuthInvalidTokenError("unknown refresh token")
	}

	if refreshToken.ExpiresAt().Before(time.Now()) {
		return nil, nil, domain.NewAuthExpiredTokenError("refresh token expired")
	}

	if err := handler.refreshTokenRepository.DeleteByValue(ctx, command.RefreshTokenString); err != nil {
		return nil, nil, err
	}

	account, err := handler.accountRepository.FindByID(ctx, refreshToken.AccountID())
	if err != nil {
		return nil, nil, err
	}

	if account == nil {
		return nil, nil, application.NewNotFoundError("not found account with given refresh token")
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
