-- +goose Up
-- +goose StatementBegin
create table accounts (
    id uuid primary key,
    login varchar not null,
    email varchar,
    password_hash varchar not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table accounts;
-- +goose StatementEnd
