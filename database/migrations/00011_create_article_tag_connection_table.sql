-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table article_tag_connections (
    id serial not null,
	primary key (id),

    article_id serial not null references articles(id),
    tag_id serial not null references tags(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table article_tag_connections;
-- +goose StatementEnd
