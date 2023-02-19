-- Add migration script here
CREATE TABLE "user" (
    id uuid not null,
    PRIMARY KEY (id),
    email text not null unique,
    hashed_password text not null
);
