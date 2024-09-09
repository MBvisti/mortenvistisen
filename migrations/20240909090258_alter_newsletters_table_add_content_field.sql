-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
alter table newsletters drop column edition, drop column body, drop column associated_article_id, add column content text not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
alter table newsletters add column edition int, add column body json not null, add column associated_article_id uuid not null, add constraint fk_associated_article_id foreign key (associated_article_id) references posts(id), drop column content;
-- +goose StatementEnd
