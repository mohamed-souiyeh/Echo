-- name: GetAllUsers :many
SELECT * FROM "users";

-- name: CreateUser :one
INSERT INTO "users" ("username", "password") 
VALUES ($1, $2) 
returning *;

-- name: GetUserByUsername :one
SELECT * FROM "users" 
WHERE "username" = $1
LIMIT 1;