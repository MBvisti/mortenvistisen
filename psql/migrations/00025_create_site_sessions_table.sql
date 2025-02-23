-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists site_sessions (
	id uuid not null,
	primary key (id),
	created_at timestamp with time zone,
	hostname varchar(100),
	browser varchar(20),
	os varchar(20),
	device varchar(20),
	screen varchar(11),
	lang varchar(35),
	country varchar(2),
	subdivision1 varchar(20),
	subdivision2 varchar(50),
	city varchar(50),
	finger varchar(500)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists site_sessions;
-- +goose StatementEnd
