package aggregates

import (
	"strings"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
)

type Account struct {
	id           *vo.AccountID
	login        string
	email        *vo.Email
	passwordHash *vo.PasswordHash
}

func (a *Account) ID() *vo.AccountID {
	return a.id
}

func (a *Account) Login() string {
	return a.login
}

func (a *Account) Email() *vo.Email {
	return a.email
}

func (a *Account) PasswordHash() *vo.PasswordHash {
	return a.passwordHash
}

func NewAccount(
	id *vo.AccountID,
	login string,
	email *vo.Email,
	passwordHash *vo.PasswordHash,
) (*Account, error) {
	var err error

	if id == nil {
		return nil, domain.NewInvalidArgumentError("account id cannot be nil")
	}

	err = checkPasswordHash(passwordHash)
	if err != nil {
		return nil, err
	}

	err = checkLogin(login)
	if err != nil {
		return nil, err
	}

	return &Account{id: id, login: login, email: email, passwordHash: passwordHash}, nil
}

func checkPasswordHash(passwordHash *vo.PasswordHash) error {
	if passwordHash == nil {
		return domain.NewInvalidArgumentError("account password hash cannot be nil")
	}

	return nil
}

func checkLogin(login string) error {
	trimmedLogin := strings.TrimSpace(login)
	if trimmedLogin == "" {
		return domain.NewInvalidArgumentError("login cannot be empty")
	}

	return nil
}

func (a *Account) ChangeLogin(newLogin string) error {
	err := checkLogin(newLogin)
	if err != nil {
		return err
	}

	a.login = newLogin
	return nil
}

func (a *Account) ChangeEmail(newEmail *vo.Email) error {
	a.email = newEmail
	return nil
}

func (a *Account) ChangePasswordHash(newPasswordHash *vo.PasswordHash) error {
	err := checkPasswordHash(newPasswordHash)
	if err != nil {
		return err
	}

	a.passwordHash = newPasswordHash
	return nil
}
