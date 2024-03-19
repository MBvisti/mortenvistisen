-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create table if not exists subscriber_tokens (
    id uuid not null,
    primary key (id),
    created_at timestamp with time zone not null,
    hash text not null,
    expires_at timestamp with time zone not null,
    scope varchar(255) not null,
    subscriber_id uuid not null,
    constraint fk_tokens_subscriber_id foreign key (subscriber_id) references subscribers(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists subscriber_tokens;
-- +goose StatementEnd
