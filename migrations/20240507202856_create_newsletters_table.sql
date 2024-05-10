-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists newsletters (
  id uuid primary key,
  created_at timestamptz not null,
  updated_at timestamptz not null,
  title varchar not null,
  edition int,
  released bool default false,
  released_at timestamptz,
  body json not null,
  associated_article_id uuid not null,
  constraint fk_associated_article_id foreign key (associated_article_id) references posts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists newsletters;
-- +goose StatementEnd
