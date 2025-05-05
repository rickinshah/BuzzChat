-- name: InsertToken :exec
INSERT INTO tokens (Hash, user_id, expiry, scope)
VALUES ($1, $2, $3, $4);

-- name: DeleteAllTokens :exec
DELETE FROM tokens
WHERE scope = $1 AND user_id = $2;

-- name: GetUserByToken :one
SELECT u.user_pid, u.username, u.email, u.name, u.password_hash, u.bio, u.activated, u.profile_pic,
       u.created_at, u.updated_at, u.version
FROM users u
INNER JOIN tokens t ON u.user_pid = t.user_id
WHERE t.hash = $1
  AND t.scope = $2
  AND t.expiry > $3;
