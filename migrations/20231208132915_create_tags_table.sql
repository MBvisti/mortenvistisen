-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table tags (
    id serial primary key,
    name varchar(255) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table tags;
-- +goose StatementEnd
