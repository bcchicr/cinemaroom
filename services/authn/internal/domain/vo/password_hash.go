package vo

import (
	"strings"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/common"
)

type PasswordHash struct {
	value string
}

func NewPasswordHash(rawHash string) (*PasswordHash, error) {
	rawHash = strings.TrimSpace(rawHash)

	if rawHash == "" {
		return nil, domain.NewInvalidArgumentError("password hash cannot be empty")
	}

	return &PasswordHash{value: rawHash}, nil
}

func (p *PasswordHash) Value() string {
	return p.value
}

func (p *PasswordHash) String() string {
	return "";
}

func (p *PasswordHash) ToMap() map[string]any {
	return map[string]any{
		"value": "",
	}
}

func (p *PasswordHash) Equals(v common.ValueObject) bool {
	other, ok := v.(*PasswordHash)
	if !ok {
		return false
	}

	return p.Value() == other.Value()
}
