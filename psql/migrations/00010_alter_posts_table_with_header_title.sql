-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
alter table posts add column header_title varchar(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
alter table posts drop column header_title;
-- +goose StatementEnd
