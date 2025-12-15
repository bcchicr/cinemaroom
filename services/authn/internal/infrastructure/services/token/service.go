package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/repositories"
	tokenService "github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/token"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type service struct {
	accessSecret           []byte
	accessTTL              time.Duration
	refreshTTL             time.Duration
	refreshTokenRepository repositories.RefreshTokenRepository
	redisClient            *redis.Client
}

func NewService(
	accessSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	refreshTokenRepository repositories.RefreshTokenRepository,
	redisClient *redis.Client,
) tokenService.Service {
	return &service{
		accessSecret:           []byte(accessSecret),
		accessTTL:              accessTTL,
		refreshTTL:             refreshTTL,
		refreshTokenRepository: refreshTokenRepository,
		redisClient:            redisClient,
	}
}

func (s *service) Generate(a *aggregates.Account) (*vo.AccessToken, *aggregates.RefreshToken, error) {
	accessToken, err := s.generateAccessToken(
		a.ID(),
		s.accessTTL,
		s.accessSecret,
	)
	if err != nil {
		var dErr domain.DomainError
		if errors.As(err, &dErr) {
			return nil, nil, err
		}
		return nil, nil, domain.NewInternalError("failed to generate access token: " + err.Error())
	}

	refreshToken, err := s.generateRefreshToken(
		a.ID(),
		s.refreshTTL,
	)
	if err != nil {
		var dErr domain.DomainError
		if errors.As(err, &dErr) {
			return nil, nil, err
		}
		return nil, nil, domain.NewInternalError("failed to generate refresh token: " + err.Error())
	}

	return accessToken, refreshToken, nil
}

func (s *service) GetClaims(ctx context.Context, rawAccessToken string) (*tokenService.Claims, error) {
	token, err := jwt.ParseWithClaims(
		rawAccessToken,
		&tokenService.Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.NewAuthInvalidTokenError(
					fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]),
				)
			}
			return s.accessSecret, nil
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired),
			errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, domain.NewAuthExpiredTokenError("token is expired or not yet valid")

		case errors.Is(err, jwt.ErrTokenMalformed),
			errors.Is(err, jwt.ErrTokenUnverifiable),
			errors.Is(err, jwt.ErrSignatureInvalid):
			return nil, domain.NewAuthInvalidTokenError("invalid token")

		default:
			return nil, domain.NewInternalError("failed to parse token: " + err.Error())
		}
	}

	claims, ok := token.Claims.(*tokenService.Claims)
	if !ok || !token.Valid {
		return nil, domain.NewAuthInvalidTokenError("invalid token")
	}

	jti := claims.ID
	isTokenRevoked, err := s.redisClient.Exists(ctx, "jwt:blacklist:"+jti).Result()
	if err != nil {
		return nil, domain.NewInternalError(
			fmt.Sprintf("redis error: %v", err),
		)
	}

	if isTokenRevoked == 1 {
		return nil, domain.NewAuthExpiredTokenError("token is revoked")
	}

	return claims, nil
}

func (s *service) generateRefreshToken(accountID *vo.AccountID, ttl time.Duration) (*aggregates.RefreshToken, error) {
	expiresAt := time.Now().Add(ttl)

	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return nil, domain.NewInternalError("failed to generate random bytes for refresh token: " + err.Error())
	}

	refreshTokenValue := base64.URLEncoding.EncodeToString(randomBytes)
	refreshTokenHash := s.Hash(refreshTokenValue)

	refreshTokenID, err := s.refreshTokenRepository.NextIdentity()
	if err != nil {
		return nil, domain.NewInternalError("failed to get next refresh token identity: " + err.Error())
	}

	return aggregates.NewRefreshToken(
		refreshTokenID,
		accountID,
		&refreshTokenValue,
		refreshTokenHash,
		expiresAt,
	)
}

func (s *service) generateAccessToken(
	accountID *vo.AccountID,
	ttl time.Duration,
	secret []byte,
) (*vo.AccessToken, error) {
	expiresAt := time.Now().Add(ttl)

	jti := uuid.NewString()
	claims := tokenService.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   accountID.String(),
		},
		AccountID: accountID.String(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	value, err := token.SignedString(secret)
	if err != nil {
		return nil, domain.NewInternalError("jwt signing failed" + err.Error())
	}

	return vo.NewAccessToken(jti, value, expiresAt)
}

func (s *service) Revoke(ctx context.Context, token *vo.AccessToken) error {
	jti := token.Jti()

	exp := token.ExpiresAt()
	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil
	}

	if err := s.redisClient.SetEx(ctx, "jwt:blacklist:"+jti, "1", ttl).Err(); err != nil {
		return domain.NewInternalError("failed to blacklist token: " + err.Error())
	}

	return nil
}

func (s *service) Hash(refreshTokenString string) string {
	hash := sha256.Sum256([]byte(refreshTokenString))
	return hex.EncodeToString(hash[:])
}
