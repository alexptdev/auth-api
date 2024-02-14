-- +goose Up
-- +goose StatementBegin

create table users (
    user_id         serial primary key,
    user_name       varchar(35),
    user_email      varchar(128),
    user_password   varchar(128),
    user_role       int,
    user_created_at timestamp not null default now(),
    user_updated_at timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists users cascade;

-- +goose StatementEnd
