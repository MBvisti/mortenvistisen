-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE articles 
RENAME COLUMN published_at TO first_published_at;

ALTER TABLE articles 
ADD COLUMN is_published boolean default false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE articles
RENAME COLUMN first_published_at TO published_at;

ALTER TABLE articles
DROP COLUMN is_published;
-- +goose StatementEnd
