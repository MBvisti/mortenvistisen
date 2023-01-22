-- Add migration script here
CREATE TABLE subscriptions (
    id uuid not null,
    PRIMARY KEY (id),
    email text not null unique,
    subscribed_at timestamptz not null,
    is_verified bool not null
);
