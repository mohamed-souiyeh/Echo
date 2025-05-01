-- name: GetAllUsers :many
SELECT * FROM "users";

-- name: CreateUser :one
INSERT INTO "users" ("username", "password") 
VALUES (?, ?) 
returning *;

-- name: GetUserByUsername :one
SELECT * FROM "users" 
WHERE "username" = ?
LIMIT 1;