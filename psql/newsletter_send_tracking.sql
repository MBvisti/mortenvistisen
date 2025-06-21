-- name: InsertNewsletterEmailSend :one
INSERT INTO newsletter_email_sends (
    id, newsletter_id, subscriber_id, status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, NOW(), NOW()
) RETURNING *;

-- name: UpdateNewsletterEmailSendStatus :one
UPDATE newsletter_email_sends 
SET 
    status = $3,
    sent_at = $4,
    failed_at = $5,
    error_message = $6,
    updated_at = NOW()
WHERE newsletter_id = $1 AND subscriber_id = $2
RETURNING *;

-- name: GetNewsletterSendStats :one
SELECT 
    newsletter_id,
    COUNT(*) as total_emails,
    COUNT(*) FILTER (WHERE status = 'sent') as sent_emails,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_emails,
    COUNT(*) FILTER (WHERE status = 'bounced') as bounced_emails,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_emails,
    ROUND(
        (COUNT(*) FILTER (WHERE status = 'sent')::DECIMAL / COUNT(*)) * 100, 2
    ) as completion_rate
FROM newsletter_email_sends 
WHERE newsletter_id = $1
GROUP BY newsletter_id;

-- name: GetAllNewsletterSendStats :many
SELECT 
    newsletter_id,
    COUNT(*) as total_emails,
    COUNT(*) FILTER (WHERE status = 'sent') as sent_emails,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_emails,
    COUNT(*) FILTER (WHERE status = 'bounced') as bounced_emails,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_emails,
    ROUND(
        (COUNT(*) FILTER (WHERE status = 'sent')::DECIMAL / COUNT(*)) * 100, 2
    ) as completion_rate
FROM newsletter_email_sends 
GROUP BY newsletter_id;

-- name: GetNewsletterEmailSendsByNewsletter :many
SELECT * FROM newsletter_email_sends 
WHERE newsletter_id = $1
ORDER BY created_at DESC;

-- name: DeleteNewsletterEmailSends :exec
DELETE FROM newsletter_email_sends 
WHERE newsletter_id = $1;
