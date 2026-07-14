-- +goose Up
CREATE UNIQUE INDEX newsletters_slug_unique
ON newsletters (slug)
WHERE slug IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS newsletters_slug_unique;
