-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table posts_tags (
    id uuid primary key,
    post_id uuid not null references posts(id),
    tag_id uuid not null references tags(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table posts_tags;
-- +goose StatementEnd
