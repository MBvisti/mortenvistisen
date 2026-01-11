-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists subscribers (
    id uuid not null,
    primary key (id),

    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,

    email varchar,
    subscribed_at timestamp with time zone,
    referer varchar,
    is_verified bool default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists subscribers;
-- +goose StatementEnd
