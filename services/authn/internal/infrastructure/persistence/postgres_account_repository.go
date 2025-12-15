package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure"
	"github.com/google/uuid"
)

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) NextIdentity() (*vo.AccountID, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, infrastructure.NewInternalError(
			fmt.Sprintf("failed to get next account ID: %v", err),
		)
	}

	return vo.NewAccountIDFromString(uuid.String())
}

func (r *PostgresAccountRepository) Save(ctx context.Context, a *aggregates.Account) error {

	query := `
        INSERT INTO accounts (id, login, email, password_hash)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET
            login = EXCLUDED.login,
            email = EXCLUDED.email,
            password_hash = EXCLUDED.password_hash,
            updated_at = current_timestamp;
    `

	emailParam := sql.NullString{}
	if a.Email() != nil {
		emailParam = sql.NullString{
			String: a.Email().String(),
			Valid:  true,
		}
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		a.ID().String(),
		a.Login(),
		emailParam,
		a.PasswordHash().Value(),
	)
	if err != nil {
		return infrastructure.NewDBTransactionFailedError(
			fmt.Sprintf("failed to save account: %v", err),
		)
	}

	return nil
}

func (r *PostgresAccountRepository) FindByID(ctx context.Context, id *vo.AccountID) (*aggregates.Account, error) {
	query := `
		SELECT id, login, email, password_hash
		FROM accounts
		WHERE id = $1
		LIMIT 1;
	`

	row := r.db.QueryRowContext(ctx, query, id.String())

	var (
		rawID           string
		rawLogin        string
		rawEmail        sql.NullString
		rawPasswordHash string
	)

	err := row.Scan(&rawID, &rawLogin, &rawEmail, &rawPasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, infrastructure.NewDBQueryFailedError(
			fmt.Sprintf("failed to scan account with id %s: %v", id.String(), err),
		)
	}

	return mapRawAccountToAggregate(rawID, rawLogin, rawEmail, rawPasswordHash)
}

func (r *PostgresAccountRepository) FindByEmail(ctx context.Context, email *vo.Email) (*aggregates.Account, error) {
	query := `
		SELECT id, login, email, password_hash
		FROM accounts
		WHERE email = $1
		LIMIT 1;
	`

	row := r.db.QueryRowContext(ctx, query, email.String())

	var (
		rawID           string
		rawLogin        string
		rawEmail        sql.NullString
		rawPasswordHash string
	)

	err := row.Scan(&rawID, &rawLogin, &rawEmail, &rawPasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, infrastructure.NewDBQueryFailedError(
			fmt.Sprintf("failed to scan account by email: %v", err),
		)
	}

	return mapRawAccountToAggregate(rawID, rawLogin, rawEmail, rawPasswordHash)
}

func (r *PostgresAccountRepository) FindByLogin(ctx context.Context, login string) (*aggregates.Account, error) {
	query := `
		SELECT id, login, email, password_hash
		FROM accounts
		WHERE login = $1
		LIMIT 1;
	`

	row := r.db.QueryRowContext(ctx, query, login)

	var (
		rawID           string
		rawLogin        string
		rawEmail        sql.NullString
		rawPasswordHash string
	)

	err := row.Scan(&rawID, &rawLogin, &rawEmail, &rawPasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, infrastructure.NewDBQueryFailedError(
			fmt.Sprintf("failed to scan account by login: %v", err),
		)
	}

	return mapRawAccountToAggregate(rawID, rawLogin, rawEmail, rawPasswordHash)
}

func mapRawAccountToAggregate(
	rawID string,
	rawLogin string,
	rawEmail sql.NullString,
	rawPasswordHash string,
) (*aggregates.Account, error) {
	id, err := vo.NewAccountIDFromString(rawID)
	if err != nil {
		return nil, err
	}

	var email *vo.Email
	if rawEmail.Valid {
		email, err = vo.NewEmail(rawEmail.String)
		if err != nil {
			return nil, err
		}
	}

	passwordHash, err := vo.NewPasswordHash(rawPasswordHash)
	if err != nil {
		return nil, err
	}

	return aggregates.NewAccount(
		id,
		rawLogin,
		email,
		passwordHash,
	)
}
