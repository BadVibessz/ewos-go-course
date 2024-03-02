-- +goose Up
-- +goose StatementBegin
CREATE TABLE private_message
(
    id        bigserial                     not null,
    from_id   integer references users (id) not null,
    to_id     integer references users (id) not null,
    content   text                          not null,
    sent_at   timestamp                     not null,
    edited_at timestamp                     not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE private_message
-- +goose StatementEnd
