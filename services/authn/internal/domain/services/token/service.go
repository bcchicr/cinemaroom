package token

import (
	"context"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	AccountID string `json:"account_id"`
}

type Service interface {
	Generate(a *aggregates.Account) (*vo.AccessToken, *aggregates.RefreshToken, error)
	GetClaims(ctx context.Context, JwtTokenString string) (*Claims, error)
	Revoke(ctx context.Context, token *vo.AccessToken) error
	Hash(refreshTokenValue string) string
}
