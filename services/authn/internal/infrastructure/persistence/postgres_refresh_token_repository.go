package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/aggregates"
	"github.com/bcchicr/cinemaroom/services/authn/internal/domain/vo"
	"github.com/bcchicr/cinemaroom/services/authn/internal/infrastructure"
	"github.com/google/uuid"
)

type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

func NewPostgresRefreshTokenRepository(db *sql.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (r *PostgresRefreshTokenRepository) NextIdentity() (*vo.RefreshTokenID, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return vo.NewRefreshTokenIdFromString(uuid.String())
}

func (r *PostgresRefreshTokenRepository) Save(ctx context.Context, a *aggregates.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (id, account_id, value_hash, expires_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET
            account_id = EXCLUDED.account_id,
            value_hash = EXCLUDED.value_hash,
            expires_at = EXCLUDED.expires_at,
            updated_at = current_timestamp;
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		a.ID().String(),
		a.AccountID().String(),
		a.ValueHash(),
		a.ExpiresAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save refresh_token: %w", err)
	}

	return nil
}

func (r *PostgresRefreshTokenRepository) FindByID(ctx context.Context, id *vo.RefreshTokenID) (*aggregates.RefreshToken, error) {
	query := `
		SELECT id, account_id, value_hash, expires_at
		FROM refresh_tokens
		WHERE id = $1
		LIMIT 1;
	`

	row := r.db.QueryRowContext(ctx, query, id.String())

	var (
		rawID        string
		rawAccountID string
		rawValueHash string
		rawExpiresAt time.Time
	)

	err := row.Scan(&rawID, &rawAccountID, &rawValueHash, &rawExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, infrastructure.NewDBQueryFailedError(
			fmt.Sprintf("failed to scan refresh token with id %s :%v", id.String(), err),
		)
	}

	return mapRawRefreshTokenToAggregate(rawID, rawAccountID, nil, rawValueHash, rawExpiresAt)
}

func (r *PostgresRefreshTokenRepository) FindByHash(ctx context.Context, valueHash string) (*aggregates.RefreshToken, error) {
	query := `
		SELECT id, account_id, value_hash, expires_at
		FROM refresh_tokens
		WHERE value_hash = $1
		LIMIT 1;
	`

	row := r.db.QueryRowContext(ctx, query, valueHash)

	var (
		rawID        string
		rawAccountID string
		rawValueHash string
		rawExpiresAt time.Time
	)

	err := row.Scan(&rawID, &rawAccountID, &rawValueHash, &rawExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, infrastructure.NewDBQueryFailedError(
			fmt.Sprintf("failed to scan refresh token by value: %v", err),
		)
	}

	return mapRawRefreshTokenToAggregate(rawID, rawAccountID, nil, rawValueHash, rawExpiresAt)
}

func mapRawRefreshTokenToAggregate(
	rawID string,
	rawAccountID string,
	rawValue *string,
	rawValueHash string,
	rawExpiresAt time.Time,
) (*aggregates.RefreshToken, error) {
	id, err := vo.NewRefreshTokenIdFromString(rawID)
	if err != nil {
		return nil, err
	}

	accountID, err := vo.NewAccountIDFromString(rawAccountID)
	if err != nil {
		return nil, err
	}

	return aggregates.NewRefreshToken(
		id,
		accountID,
		rawValue,
		rawValueHash,
		rawExpiresAt,
	)
}

func (r *PostgresRefreshTokenRepository) Delete(ctx context.Context, a *aggregates.RefreshToken) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1;
	`

	id := a.ID().String()

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return infrastructure.NewDBTransactionFailedError(
			fmt.Sprintf("failed to delete refresh token with id %s: %v", id, err),
		)
	}

	return nil
}
