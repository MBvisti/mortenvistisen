-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS newsletter_email_sends (
    id uuid not null,
    primary key (id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    status VARCHAR NOT NULL,

    sent_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    error_message TEXT,

	newsletter_id UUID NOT NULL REFERENCES newsletters(id) ON DELETE CASCADE,
	subscriber_id UUID NOT NULL REFERENCES subscribers(id) ON DELETE CASCADE,
    UNIQUE(newsletter_id, subscriber_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS newsletter_email_sends;
-- +goose StatementEnd
