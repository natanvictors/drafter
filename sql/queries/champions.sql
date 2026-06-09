-- name: AddChampions :one

INSERT INTO champions (champ_name)
VALUES ($1)
ON CONFLICT (champ_name) DO UPDATE SET
    champ_name = EXCLUDED.champ_name
RETURNING *;

