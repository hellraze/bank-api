-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA bank;
CREATE TABLE bank.user (
                           user_id uuid PRIMARY KEY,
                           login text UNIQUE NOT NULL,
                           hash_password bytea NOT NULL
);
CREATE TABLE bank.account (
                              account_id uuid PRIMARY KEY,
                              name text UNIQUE NOT NULL,
                              balance int,
                              user_id uuid REFERENCES bank.user(user_id) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bank.account;
DROP TABLE bank.user;
DROP SCHEMA bank;
-- +goose StatementEnd
