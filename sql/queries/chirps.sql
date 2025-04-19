-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirpsAsc :many
SELECT
    id,
    created_at,
    updated_at,
    body,
    user_id
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsById :one
SELECT * FROM chirps
WHERE id = $1;
