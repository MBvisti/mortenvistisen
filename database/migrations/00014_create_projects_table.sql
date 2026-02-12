-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists projects (
    id serial not null,
    primary key (id),

    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,

    published boolean not null default false,

    title varchar(120) not null,
    slug varchar(255) not null unique,
    started_at timestamp with time zone,
    status varchar(80) not null default 'planned',
    description text not null default '',
    content text not null default '',
    project_url varchar(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists projects;
-- +goose StatementEnd
