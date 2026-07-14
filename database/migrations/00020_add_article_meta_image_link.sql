-- +goose Up
-- +goose StatementBegin
ALTER TABLE articles ADD COLUMN meta_image_link text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE articles DROP COLUMN meta_image_link;
-- +goose StatementEnd
