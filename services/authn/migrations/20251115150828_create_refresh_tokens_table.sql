-- +goose Up
-- +goose StatementBegin
create table refresh_tokens (
    id uuid primary key,
    account_id uuid references accounts(id),
    value varchar not null,
    expires_at timestamp not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table refresh_tokens;
-- +goose StatementEnd
