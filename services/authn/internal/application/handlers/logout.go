package handlers

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/application"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/repositories"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/token"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/redis/go-redis/v9"
)

type LogoutCommand struct {
	JwtTokenString     string
	RefreshTokenString string
}

type LogoutHandler interface {
	Handle(context.Context, LogoutCommand) (*vo.AccountID, error)
}

type logoutHandler struct {
	tokenService           token.Service
	accountRepository      repositories.AccountRepository
	refreshTokenRepository repositories.RefreshTokenRepository
	redisClient            *redis.Client
}

func NewLogoutHandler(
	tokenService token.Service,
	accountRepository repositories.AccountRepository,
	refreshTokenRepository repositories.RefreshTokenRepository,
	redisClient *redis.Client,
) LogoutHandler {
	return &logoutHandler{
		tokenService:           tokenService,
		accountRepository:      accountRepository,
		refreshTokenRepository: refreshTokenRepository,
		redisClient:            redisClient,
	}
}

func (handler *logoutHandler) Handle(ctx context.Context, command LogoutCommand) (*vo.AccountID, error) {
	claims, err := handler.tokenService.GetClaims(ctx, command.JwtTokenString)
	if err != nil {
		return nil, err
	}

	token, err := vo.NewAccessToken(
		claims.ID,
		command.JwtTokenString,
		claims.ExpiresAt.Time,
	)
	if err != nil {
		return nil, err
	}

	if err := handler.tokenService.Revoke(ctx, token); err != nil {
		return nil, err
	}

	accountID, err := vo.NewAccountIDFromString(claims.AccountID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := handler.refreshTokenRepository.FindByValue(ctx, command.RefreshTokenString)
	if err != nil {
		return nil, err
	}

	if !refreshToken.AccountID().Equals(accountID) {
		return nil, application.NewNotAuthorizedError("user is not authorized to delete other user's refresh token")
	}

	if err := handler.refreshTokenRepository.DeleteByValue(ctx, command.RefreshTokenString); err != nil {
		return nil, err
	}

	return accountID, nil
}
