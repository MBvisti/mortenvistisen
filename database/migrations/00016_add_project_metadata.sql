-- +goose Up
-- +goose StatementBegin
alter table projects
    add column if not exists meta_title varchar(60) not null default '',
    add column if not exists meta_description varchar(160) not null default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table projects
    drop column if exists meta_description,
    drop column if exists meta_title;
-- +goose StatementEnd
