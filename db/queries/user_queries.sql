-- name: GetAllUsers :many
SELECT * FROM "users";

-- name: CreateUser :one
INSERT INTO "users" ("username", "password") 
VALUES (?, ?) 
returning *;