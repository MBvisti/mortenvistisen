-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists articles (
    id serial not null,
    primary key (id),

    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    first_published_at timestamp with time zone,

	published boolean not null default false,

	title varchar(100) not null,
	excerpt varchar(255),
	meta_title varchar(100),
	meta_description varchar(160),
	slug varchar(255) not null unique,
	image_link varchar(255),
	read_time integer,

	content text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists articles;
-- +goose StatementEnd
