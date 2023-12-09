-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table posts (
    id serial primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title varchar(255) not null,
    filename varchar(255) not null,
    slug varchar(255) not null,
    excerpt text not null,
    draft boolean not null default true,
    released_at timestamp,
    read_time integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table posts;
-- +goose StatementEnd
