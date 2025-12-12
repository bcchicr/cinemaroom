-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    account_id UUID REFERENCES accounts(id),
    value_hash VARCHAR NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX refresh_tokens_value_hash_idx ON refresh_tokens(value_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
