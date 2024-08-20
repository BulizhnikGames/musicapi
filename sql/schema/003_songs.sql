-- +goose Up

CREATE TABLE songs(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    UNIQUE (name, album_id)
);

-- +goose Down

DROP TABLE songs;