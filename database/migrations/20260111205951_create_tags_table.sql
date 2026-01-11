-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists article_tags (
    id uuid not null,
    primary key (id),

    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,

	title varchar(255) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists article_tags;
-- +goose StatementEnd
