-- +goose Up
-- +goose StatementBegin
ALTER TABLE newsletters ADD COLUMN meta_image_link text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE newsletters DROP COLUMN meta_image_link;
-- +goose StatementEnd
