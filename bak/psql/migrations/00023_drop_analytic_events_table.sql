-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE IF EXISTS analytic_events;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
create table if not exists analytic_events (
    id UUID NOT NULL PRIMARY KEY,
	analytics_id uuid,
    element_tag varchar,
	element_text varchar,
	element_class varchar,
	element_id varchar,
	element_href varchar,
	custom_data varchar
);
-- +goose StatementEnd
