-- +goose Up

CREATE TABLE playlists_songs(
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    song_id UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    add_time TIMESTAMP NOT NULL,
    UNIQUE (playlist_id, song_id)
);

-- +goose Down

DROP TABLE playlists_songs;