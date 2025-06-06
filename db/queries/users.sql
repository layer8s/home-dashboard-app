-- name: UpsertUser :one
INSERT INTO users (
    name,
    email,
    activated,
    provider,
    provider_id
)
VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (provider, provider_id) DO UPDATE
SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    activated = EXCLUDED.activated,
    version = users.version + 1,
    created_at = NOW()
RETURNING id, created_at, version;

-- name: GetUserByProvider :one
SELECT id, name
FROM users 
WHERE provider = $1 AND provider_id = $2;

-- -- name: GetUserByEmail :one
-- SELECT id, created_at, name, email, password_hash, activated, version
-- FROM users
-- WHERE email = $1;

-- -- name: UpdateUser :one
-- UPDATE users 
-- SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
-- WHERE id = $5 AND version = $6
-- RETURNING version;