-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
alter table users drop column name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
alter table users add column name varchar(255) not null;
-- +goose StatementEnd
