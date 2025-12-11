package password

import (
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type Service interface {
	Hash(p *vo.Password) (*vo.PasswordHash, error)
	Verify(p *vo.Password, h *vo.PasswordHash) bool
}
