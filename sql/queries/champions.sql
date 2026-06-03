-- name: AddChampions :one

INSERT INTO champions (champ_name)
VALUES ($1)
RETURNING *;

-- name: UpdateStatus :one
UPDATE champions
SET
    times_picked=$1,
    picked_top=$2,
    picked_jungle=$3,
    picked_mid=$4,
    picked_adc=$5,
    picked_sup=$6,
    wins_top=$7,
    losses_top=$8,
    wins_jg=$9,
    losses_jg=$10,
    wins_mid=$11,
    losses_mid=$12,
    wins_adc=$13,
    losses_adc=$14,
    wins_sup=$15,
    losses_sup=$16
WHERE champ_name = $17
RETURNING *;

-- name: ChangeRoles :one
UPDATE champions
SET ingame_role = $1
WHERE champ_name = $2
RETURNING *;