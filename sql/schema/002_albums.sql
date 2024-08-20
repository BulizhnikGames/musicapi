-- +goose Up

CREATE TABLE albums(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    artist_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (name, artist_id)
);

-- +goose Down

DROP TABLE albums;