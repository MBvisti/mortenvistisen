-- +goose Up
-- +goose StatementBegin
ALTER TABLE newsletters ADD COLUMN image_link text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE newsletters DROP COLUMN image_link;
-- +goose StatementEnd
