package vo

import (
	"strings"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/common"
)

type AccessToken struct {
	jti       string
	value     string
	expiresAt time.Time
}

func NewAccessToken(jti string, value string, expiresAt time.Time) (*AccessToken, error) {
	trimmedJti := strings.TrimSpace(jti)
	if trimmedJti == "" {
		return nil, domain.NewInvalidArgumentError("token jti cannot be empty")
	}

	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return nil, domain.NewInvalidArgumentError("token value cannot be empty")
	}

	return &AccessToken{jti: jti, value: trimmedValue, expiresAt: expiresAt}, nil
}

func (p *AccessToken) Jti() string {
	return p.jti
}

func (p *AccessToken) Value() string {
	return p.value
}

func (p *AccessToken) String() string {
	return p.Value()
}

func (p *AccessToken) ExpiresAt() time.Time {
	return p.expiresAt
}

func (p *AccessToken) ToMap() map[string]any {
	return map[string]any{
		"jti":       p.Jti(),
		"value":     p.Value(),
		"expiresAt": p.ExpiresAt(),
	}
}

func (p *AccessToken) Equals(v common.ValueObject) bool {
	other, ok := v.(*AccessToken)
	if !ok {
		return false
	}

	return p.Jti() == other.Jti() &&
		p.Value() == other.Value() &&
		p.ExpiresAt().Equal(other.ExpiresAt())
}
