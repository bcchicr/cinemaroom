package aggregates

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type RefreshToken struct {
	id        *vo.RefreshTokenID
	accountID *vo.AccountID
	value     string
	expiresAt time.Time
}

func (a *RefreshToken) ID() *vo.RefreshTokenID {
	return a.id
}

func (a *RefreshToken) AccountID() *vo.AccountID {
	return a.accountID
}

func (a *RefreshToken) Value() string {
	return a.value
}

func (a *RefreshToken) Hash() (string, error) {
	hash := sha256.Sum256([]byte(a.Value()))
	return hex.EncodeToString(hash[:]), nil
}

func (a *RefreshToken) ExpiresAt() time.Time {
	return a.expiresAt
}

func NewRefreshToken(
	id *vo.RefreshTokenID,
	accountID *vo.AccountID,
	value string,
	expiresAt time.Time,
) (*RefreshToken, error) {
	if id == nil {
		return nil, domain.NewInvalidArgumentError("refresh token id cannot be nil")
	}

	if accountID == nil {
		return nil, domain.NewInvalidArgumentError("refresh token account id cannot be nil")
	}

	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return nil, domain.NewInvalidArgumentError("token value cannot be empty")
	}

	return &RefreshToken{id: id, accountID: accountID, value: value, expiresAt: expiresAt}, nil
}
