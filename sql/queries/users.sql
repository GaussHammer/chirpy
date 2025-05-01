-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: SelectUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE email = $1;

-- name: UpdateUserById :exec
UPDATE users
SET email = $1, hashed_password = $2
WHERE id = $3;

-- name: SelectUserById :one
SELECT users.created_at, users.updated_at, refresh_tokens.token FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE users.id = $1;

-- name: UpgradeChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1;