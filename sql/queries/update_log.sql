-- name: UpsertRecord :exec

INSERT INTO update_log (
    region,
    last_updated_at
)VALUES(
    $1, $2
)
ON CONFLICT (region) DO UPDATE SET
    last_updated_at = EXCLUDED.last_updated_at;


-- name: GetUpdateRecord :one
SELECT * from update_log WHERE region = $1;