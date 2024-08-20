-- +goose Up

CREATE TABLE likes(
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    song_id UUID NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    UNIQUE (user_id, song_id)
);

-- +goose Down

DROP TABLE likes;