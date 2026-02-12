-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists newsletters (
  id serial not null,
  primary key (id),

  created_at timestamp with time zone not null,
  updated_at timestamp with time zone not null,

  title varchar(100) not null,
  slug varchar,

  meta_title varchar(100) not null,
  meta_description varchar(160) not null,

  is_published bool default false,
  released_at timestamptz,

  content text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists newsletters;
-- +goose StatementEnd
