-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS newsletter_email_sends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    newsletter_id UUID NOT NULL REFERENCES newsletters(id) ON DELETE CASCADE,
    subscriber_id UUID NOT NULL REFERENCES subscribers(id) ON DELETE CASCADE,
    email_address VARCHAR NOT NULL,
    status VARCHAR NOT NULL CHECK (status IN ('pending', 'sent', 'failed', 'bounced')),
    sent_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    error_message TEXT,
    river_job_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(newsletter_id, subscriber_id)
);

CREATE INDEX IF NOT EXISTS idx_newsletter_email_sends_newsletter_id ON newsletter_email_sends(newsletter_id);
CREATE INDEX IF NOT EXISTS idx_newsletter_email_sends_status ON newsletter_email_sends(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS newsletter_email_sends;
-- +goose StatementEnd