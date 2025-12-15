package services

import (
	"context"
	"testing"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/app"
	"github.com/bcchicr/cinemaroom/services/authn/config"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure/services/token"
	"github.com/google/uuid"
)

func TestTokenService(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("cannot get config")
	}

	container, err := app.NewContainer(cfg)
	if err != nil {
		t.Fatalf("Failed to build container: %v", err)
	}
	defer container.Close()

	accountID, _ := vo.NewAccountIDFromString(uuid.NewString())
	login := "login"
	email, _ := vo.NewEmail("email@example.com")
	passwordHash, _ := vo.NewPasswordHash("hash")

	account, _ := aggregates.NewAccount(
		accountID,
		login,
		email,
		passwordHash,
	)

	refreshTokenRepo := container.RefreshTokenRepository
	redis := container.Redis

	s := token.NewService(
		"secret",
		time.Minute,
		time.Hour,
		refreshTokenRepo,
		redis,
	)

	t.Run("generate tokens", func(t *testing.T) {
		access, refresh, err := s.Generate(account)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if access == nil {
			t.Fatalf("expected access token")
		}

		if refresh == nil {
			t.Fatalf("expected refresh token")
		}
	})

	t.Run("get claims from valid token", func(t *testing.T) {
		access, _, _ := s.Generate(account)

		claims, err := s.GetClaims(context.Background(), access.Value())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if claims.AccountID != accountID.String() {
			t.Fatalf("unexpected account id in claims")
		}
	})

	t.Run("revoke token", func(t *testing.T) {
		access, _, _ := s.Generate(account)

		err := s.Revoke(context.Background(), access)
		if err != nil {
			t.Fatalf("unexpected revoke error: %v", err)
		}

		_, err = s.GetClaims(context.Background(), access.Value())
		if err == nil {
			t.Fatalf("expected revoked token error")
		}
	})

	t.Run("hash refresh token", func(t *testing.T) {
		value := "refresh-token-value"

		hash1 := s.Hash(value)
		hash2 := s.Hash(value)

		if hash1 != hash2 {
			t.Fatalf("expected deterministic hash")
		}
	})
}
