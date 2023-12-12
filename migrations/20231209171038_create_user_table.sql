-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA bank;
CREATE TABLE bank.user (
    user_id uuid PRIMARY KEY,
    login text UNIQUE NOT NULL,
    hash_password bytea NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bank.user;
DROP SCHEMA bank;
-- +goose StatementEnd
