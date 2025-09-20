-- name: CreateUser :exec
INSERT INTO users (
    id, email, username
) VALUES (
    ?, ?, ?
);

-- name: GetUsers :many
SELECT * FROM users LIMIT ?;

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;
