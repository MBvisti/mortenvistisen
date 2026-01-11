-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table article_tag_connections (
    id uuid primary key,
    article_id uuid not null references articles(id),
    tag_id uuid not null references article_tags(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table article_tag_connections;
-- +goose StatementEnd
