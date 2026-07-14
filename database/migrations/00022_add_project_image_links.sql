-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects
    ADD COLUMN image_link text,
    ADD COLUMN meta_image_link text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects
    DROP COLUMN meta_image_link,
    DROP COLUMN image_link;
-- +goose StatementEnd
