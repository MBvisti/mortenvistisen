-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table posts_tags (
    id serial primary key,
    post_id integer not null references posts(id),
    tag_id integer not null references tags(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table posts_tags;
-- +goose StatementEnd
