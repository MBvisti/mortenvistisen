-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists analytics (
    id UUID NOT NULL PRIMARY KEY,
	timestamp TIMESTAMP WITH TIME ZONE,
    website_id UUID,
    type VARCHAR,
    url VARCHAR,
    path VARCHAR,
    referrer VARCHAR,  
    title VARCHAR,
    screen VARCHAR,
    language VARCHAR,
    visitor_id UUID,
    session_id UUID,
    scroll_depth integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists analytics;
-- +goose StatementEnd
