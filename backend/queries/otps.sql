-- name: InsertOTP :exec
INSERT INTO otps (user_pid, created_at, secret_key, expiry)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_pid)
DO UPDATE
SET created_at = $2, secret_key = $3, expiry = $4;

-- name: DeleteOTP :exec
DELETE FROM otps WHERE user_pid =$1;

-- name: GetOTP :one
SELECT user_pid, created_at, secret_key, expiry
FROM otps WHERE user_pid = $1;
