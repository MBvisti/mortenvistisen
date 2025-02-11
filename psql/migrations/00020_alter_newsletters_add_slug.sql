-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
alter table newsletters add column slug varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
alter table newsletters drop column slug;
-- +goose StatementEnd
