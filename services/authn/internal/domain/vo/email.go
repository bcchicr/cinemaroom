package vo

import (
	"net/mail"
	"strings"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/common"
)

type Email struct {
	value string
}

func NewEmail(rawEmail string) (*Email, error) {
	rawEmail = strings.TrimSpace(rawEmail)
	if rawEmail == "" {
		return nil, domain.NewInvalidArgumentError("email cannot be empty")
	}

	_, err := mail.ParseAddress(rawEmail)
	if err != nil {
		return nil, domain.NewFailedInvariantError("email string must be valid")
	}

	return &Email{value: rawEmail}, nil
}

func (p *Email) Value() string {
	return p.value
}

func (p *Email) String() string {
	return p.Value()
}

func (p *Email) ToMap() map[string]any {
	return map[string]any{
		"value": p.Value(),
	}
}

func (p *Email) Equals(v common.ValueObject) bool {
	other, ok := v.(*Email)
	if !ok {
		return false
	}

	return p.Value() == other.Value()
}
