-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS subscribers_email_unique_idx
    ON subscribers (email)
    WHERE email IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS subscribers_email_unique_idx;
-- +goose StatementEnd
