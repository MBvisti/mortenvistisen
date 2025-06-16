-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists newsletters (
  id uuid primary key,
  created_at timestamptz not null,
  updated_at timestamptz not null,

  title varchar not null,
  is_published bool default false,
  released_at timestamptz,
  slug varchar,
  content text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists newsletters;
-- +goose StatementEnd
