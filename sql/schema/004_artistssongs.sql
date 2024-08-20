-- +goose Up

CREATE TABLE artists_songs(
    artist_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    song_id UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    UNIQUE (artist_id, song_id)
);

-- +goose Down

DROP TABLE artists_songs;