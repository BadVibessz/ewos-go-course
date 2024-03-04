-- +goose Up
-- +goose StatementBegin
CREATE TABLE public_message
(
    id            bigserial primary key                                      not null,
    from_username varchar(128) references users (username) on update cascade not null,
    content       text                                                       not null,
    sent_at       timestamp                                                  not null,
    edited_at     timestamp                                                  not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public_message
-- +goose StatementEnd
