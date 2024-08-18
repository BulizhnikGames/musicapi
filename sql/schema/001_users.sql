-- +goose Up

CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    UNIQUE (email, password),
    is_artist BOOLEAN NOT NULL
);

-- +goose Down

DROP TABLE users;