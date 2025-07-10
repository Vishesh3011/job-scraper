-- name: CreateUser :exec
INSERT INTO job_scraper_users (id, name, email, location, keywords, cookie, csrf_token)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByEmail :one
SELECT * FROM job_scraper_users WHERE email = ?;

-- name: CheckUserExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM job_scraper_users WHERE email = ?) AS `exists`;

-- name: UpdateUser :exec
UPDATE job_scraper_users
SET name = ?, location = ?, keywords = ?, cookie = ?, csrf_token = ?
WHERE email = ?;

-- name: GetAllUsers :many
SELECT * FROM job_scraper_users;

-- name: GetUserByID :one
SELECT id, name, email, location, keywords, cookie, csrf_token, created_at
FROM job_scraper_users
WHERE id = ?;
