-- name: CreateUser :exec
INSERT INTO job_scraper_users (name, email, location, keywords, cookie, csrf_token)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetUserByEmail :one
SELECT * FROM job_scraper_users WHERE email = ?;

-- name: CheckUserExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM job_scraper_users WHERE email = ?) AS `exists`;