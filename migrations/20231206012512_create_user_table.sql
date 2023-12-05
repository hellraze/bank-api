-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA bank;
CREATE TABLE bank.user (
    id uuid PRIMARY KEY,
    login text NOT NULL,
    hash_password text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bank.user;
DROP SCHEMA bank;
-- +goose StatementEnd
