-- +goose Up
-- +goose StatementBegin
CREATE TABLE public_message
(
    id        bigserial                     not null,
    from_id   integer references users (id) not null,
    content   text                          not null,
    sent_at   timestamp                     not null,
    edited_at timestamp                     not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public_message
-- +goose StatementEnd
