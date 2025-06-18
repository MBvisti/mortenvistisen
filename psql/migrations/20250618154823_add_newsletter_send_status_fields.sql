-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- Add send status tracking fields to newsletters table
ALTER TABLE newsletters ADD COLUMN send_status varchar(20) DEFAULT 'draft' NOT NULL;
ALTER TABLE newsletters ADD COLUMN total_recipients int DEFAULT 0 NOT NULL;
ALTER TABLE newsletters ADD COLUMN emails_sent int DEFAULT 0 NOT NULL;
ALTER TABLE newsletters ADD COLUMN sending_started_at timestamptz;
ALTER TABLE newsletters ADD COLUMN sending_completed_at timestamptz;

-- Add check constraint for send_status
ALTER TABLE newsletters ADD CONSTRAINT chk_newsletter_send_status 
    CHECK (send_status IN ('draft', 'ready_to_send', 'sending', 'sent'));

-- Add index for querying newsletters by send status
CREATE INDEX idx_newsletters_send_status ON newsletters(send_status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

-- Remove the added fields and constraints
DROP INDEX IF EXISTS idx_newsletters_send_status;
ALTER TABLE newsletters DROP CONSTRAINT IF EXISTS chk_newsletter_send_status;
ALTER TABLE newsletters DROP COLUMN IF EXISTS sending_completed_at;
ALTER TABLE newsletters DROP COLUMN IF EXISTS sending_started_at;
ALTER TABLE newsletters DROP COLUMN IF EXISTS emails_sent;
ALTER TABLE newsletters DROP COLUMN IF EXISTS total_recipients;
ALTER TABLE newsletters DROP COLUMN IF EXISTS send_status;

-- +goose StatementEnd
