-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, password, is_artist)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByEmailAndPassword :one
SELECT * FROM users WHERE email = $1 AND password = $2;