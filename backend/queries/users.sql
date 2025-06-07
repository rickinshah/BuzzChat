-- name: GetUser :one
SELECT
    user_pid,
    username,
    email,
    name,
    password_hash,
    bio,
    activated,
    profile_pic,
    created_at,
    updated_at,
    version
FROM
    users
WHERE
    user_pid = $1;

-- name: GetUserByUsername :one
SELECT
    user_pid,
    username,
    email,
    name,
    password_hash,
    bio,
    activated,
    profile_pic,
    created_at,
    updated_at,
    version
FROM
    users
WHERE
    username = $1;

-- name: GetUserByEmailOrUsername :one
SELECT
    user_pid,
    username,
    email,
    name,
    password_hash,
    bio,
    activated,
    profile_pic,
    created_at,
    updated_at,
    version
FROM
    users
WHERE
    email = $1
    OR username = $1;

-- name: GetUserByEmail :one
SELECT
    user_pid,
    username,
    email,
    name,
    password_hash,
    bio,
    activated,
    profile_pic,
    created_at,
    updated_at,
    version
FROM
    users
WHERE
    email = $1;

-- name: UpdateUser :one
UPDATE
    users
SET
    username = $3,
    email = $4,
    name = $5,
    bio = $6,
    activated = $7,
    profile_pic = $8,
    updated_at = now(),
    version = version + 1
WHERE
    user_pid = $1
    AND version = $2
RETURNING
    updated_at,
    version;

-- name: UpdatePassword :one
UPDATE
    users
SET
    password_hash = $3,
    updated_at = now(),
    version = version + 1
WHERE
    user_pid = $1
    AND version = $2
RETURNING
    updated_at,
    version;

-- name: InsertUser :one
INSERT INTO users (username, email, name, password_hash, bio, profile_pic)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    user_pid, created_at, updated_at, version;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_pid = $1;

