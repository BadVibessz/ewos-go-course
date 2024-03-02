-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id              bigserial    not null primary key,
    email           varchar(256) not null,
    username        varchar(128) not null,
    hashed_password text         not null,
    created_at      timestamp    not null,
    edited_at       timestamp    not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd
