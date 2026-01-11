-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists articles (
    id uuid not null,
    primary key (id),

    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    first_published_at timestamp with time zone,

	title varchar(100) not null,
	excerpt varchar(255) not null,
	meta_title varchar(100) not null,
	meta_description varchar(160) not null,
	slug varchar(255) not null,
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
