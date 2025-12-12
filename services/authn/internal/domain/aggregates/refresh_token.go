package aggregates

import (
	"strings"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type RefreshToken struct {
	id        *vo.RefreshTokenID
	accountID *vo.AccountID
	value     *string
	valueHash string
	expiresAt time.Time
}

func (a *RefreshToken) ID() *vo.RefreshTokenID {
	return a.id
}

func (a *RefreshToken) AccountID() *vo.AccountID {
	return a.accountID
}

func (a *RefreshToken) Value() *string {
	return a.value
}

func (a *RefreshToken) ValueHash() string {
	return a.valueHash
}

func (a *RefreshToken) ExpiresAt() time.Time {
	return a.expiresAt
}

func NewRefreshToken(
	id *vo.RefreshTokenID,
	accountID *vo.AccountID,
	value *string,
	valueHash string,
	expiresAt time.Time,
) (*RefreshToken, error) {
	if id == nil {
		return nil, domain.NewInvalidArgumentError("refresh token id cannot be nil")
	}

	if accountID == nil {
		return nil, domain.NewInvalidArgumentError("refresh token account id cannot be nil")
	}

	trimmedValueHash := strings.TrimSpace(valueHash)
	if trimmedValueHash == "" {
		return nil, domain.NewInvalidArgumentError("token value hash cannot be empty")
	}

	if value != nil {
		trimmedValue := strings.TrimSpace(*value)
		if trimmedValue == "" {
			return nil, domain.NewInvalidArgumentError("token value cannot be empty")
		}

		if trimmedValue == trimmedValueHash {
			return nil, domain.NewFailedInvariantError("token value and value hash cannot be equal")
		}
	}

	return &RefreshToken{id: id, accountID: accountID, value: value, valueHash: valueHash, expiresAt: expiresAt}, nil
}
