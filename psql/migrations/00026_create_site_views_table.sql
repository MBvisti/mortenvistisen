-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create type site_event AS ENUM('page_view', 'page_leave', 'click');

create table if not exists site_views (
	id bigserial not null,
	primary key (id),
	session_id uuid,
	visitor_id uuid,
	created_at timestamp with time zone,
	url_path varchar(500),
	url_query varchar(500),
	referrer_path varchar(500),
	referrer_query varchar(500),
	referrer_domain varchar(500),
	page_title varchar(500),
	event_type site_event,
	event_name varchar(50),
	CONSTRAINT fk_website_views_sessions FOREIGN KEY (session_id) REFERENCES site_sessions(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists site_views;
drop type if exists site_event;
-- +goose StatementEnd
