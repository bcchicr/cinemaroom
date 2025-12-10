package vo

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/common"
)

type Password struct {
	value string
}

func NewPassword(rawPassword string) (*Password, error) {
	rawPassword = strings.TrimSpace(rawPassword)
	if rawPassword == "" {
		return nil, domain.NewInvalidArgumentError("password cannot be empty")
	}

	rawPasswordLength := utf8.RuneCountInString(rawPassword)
	if rawPasswordLength < common.PasswordMinLength {
		return nil, domain.NewFailedInvariantError(
			fmt.Sprintf("password min length: %d, got %d", common.PasswordMinLength, rawPasswordLength),
		)
	}

	rawPasswordByteLength := len(rawPassword)
	if rawPasswordByteLength > common.PasswordMaxByteLength {
		return nil, domain.NewFailedInvariantError(
			fmt.Sprintf("password max byte length: %d, got %d", common.PasswordMaxByteLength, rawPasswordByteLength),
		)
	}

	return &Password{value: rawPassword}, nil
}

func (p *Password) Value() string {
	return p.value
}

func (p *Password) String() string {
	return p.Value()
}

func (p *Password) ToMap() map[string]any {
	return map[string]any{
		"value": "",
	}
}

func (p *Password) Equals(v common.ValueObject) bool {
	other, ok := v.(*Password)
	if !ok {
		return false
	}

	return p.Value() == other.Value()
}
