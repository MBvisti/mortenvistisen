-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE articles add column read_time integer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE articles drop column read_time;
-- +goose StatementEnd
