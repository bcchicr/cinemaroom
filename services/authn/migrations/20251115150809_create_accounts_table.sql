-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    login VARCHAR NOT NULL,
    email VARCHAR,
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX accounts_login_idx ON accounts(login);
CREATE UNIQUE INDEX accounts_email_idx ON accounts(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
-- +goose StatementEnd
