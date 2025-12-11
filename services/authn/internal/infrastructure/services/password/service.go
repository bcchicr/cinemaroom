package password

import (
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/services/password"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"golang.org/x/crypto/bcrypt"
)

type service struct{}

func NewService() password.Service {
	return &service{}
}

func (s *service) Hash(p *vo.Password) (*vo.PasswordHash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p.Value()), 14)
	if err != nil {
		return nil, domain.NewInternalError("unable to generate password hash: " + err.Error())
	}
	return vo.NewPasswordHash(string(bytes))
}

func (s *service) Verify(p *vo.Password, h *vo.PasswordHash) bool {
	return bcrypt.CompareHashAndPassword([]byte(h.Value()), []byte(p.Value())) == nil
}
