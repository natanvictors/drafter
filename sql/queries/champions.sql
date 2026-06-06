-- name: AddChampions :one

INSERT INTO champions (champ_name)
VALUES ($1)
ON CONFLICT (champ_name) DO UPDATE SET
    champ_name = EXCLUDED.champ_name
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

-- name: UpdateChampionStats :exec

WITH reset AS (
    UPDATE champions SET 
  times_picked = 0, picked_top = 0, picked_jungle = 0,
  picked_mid = 0, picked_adc = 0, picked_sup = 0,
  wins_top = 0, losses_top = 0, wins_jg = 0, losses_jg = 0,
  wins_mid = 0, losses_mid = 0, wins_adc = 0, losses_adc = 0,
  wins_sup = 0, losses_sup = 0
),
picks AS (
    SELECT champ, role, won FROM (
        SELECT team1_pick1 AS champ, team1_role1 AS role, (winner='1') AS won FROM games
        UNION ALL
        SELECT team1_pick2 AS champ, team1_role2 AS role, (winner='1') AS won FROM games
        UNION ALL
        SELECT team1_pick3 AS champ, team1_role3 AS role, (winner='1') AS won FROM games
        UNION ALL
        SELECT team1_pick4 AS champ, team1_role4 AS role, (winner='1') AS won FROM games
        UNION ALL
        SELECT team1_pick5 AS champ, team1_role5 AS role, (winner='1') AS won FROM games
        UNION ALL
        SELECT team2_pick1 AS champ, team2_role1 AS role, (winner='2') AS won FROM games
        UNION ALL
        SELECT team2_pick2 AS champ, team2_role2 AS role, (winner='2') AS won FROM games
        UNION ALL
        SELECT team2_pick3 AS champ, team2_role3 AS role, (winner='2') AS won FROM games
        UNION ALL
        SELECT team2_pick4 AS champ, team2_role4 AS role, (winner='2') AS won FROM games
        UNION ALL
        SELECT team2_pick5 AS champ, team2_role5 AS role, (winner='2') AS won FROM games
    ) sub), agg AS (
        SELECT
        champ,
        COUNT(*) AS times_picked,
        COUNT(*) FILTER (WHERE role = 'Top') AS picked_top,
        COUNT(*) FILTER (WHERE role = 'Jungle') AS picked_jungle,
        COUNT(*) FILTER (WHERE role = 'Mid') AS picked_mid,
        COUNT(*) FILTER (WHERE role = 'Bot') AS picked_adc,
        COUNT(*) FILTER (WHERE role = 'Support') AS picked_sup,
        COUNT(*) FILTER (WHERE role = 'Top' AND won) AS wins_top,
        COUNT(*) FILTER (WHERE role = 'Jungle' AND won) AS wins_jg,
        COUNT(*) FILTER (WHERE role = 'Mid' AND won) AS wins_mid,
        COUNT(*) FILTER (WHERE role = 'Bot' AND won) AS wins_adc,
        COUNT(*) FILTER (WHERE role = 'Support' AND won) AS wins_sup,
        COUNT(*) FILTER (WHERE role = 'Top' AND NOT won) AS losses_top,
        COUNT(*) FILTER (WHERE role = 'Jungle' AND NOT won) AS losses_jg,
        COUNT(*) FILTER (WHERE role = 'Mid' AND NOT won) AS losses_mid,
        COUNT(*) FILTER (WHERE role = 'Bot' AND NOT won) AS losses_adc,
        COUNT(*) FILTER (WHERE role = 'Support' AND NOT won) AS losses_sup
        FROM picks
        GROUP BY champ
    )
    UPDATE champions c 
    SET
        times_picked  = a.times_picked,
        picked_top    = a.picked_top,
        picked_jungle = a.picked_jungle,
        picked_mid    = a.picked_mid,
        picked_adc    = a.picked_adc,
        picked_sup    = a.picked_sup,
        wins_top      = a.wins_top,
        losses_top    = a.losses_top,
        wins_jg       = a.wins_jg,
        losses_jg     = a.losses_jg,
        wins_mid      = a.wins_mid,
        losses_mid    = a.losses_mid,
        wins_adc      = a.wins_adc,
        losses_adc    = a.losses_adc,
        wins_sup      = a.wins_sup,
        losses_sup    = a.losses_sup
    FROM agg a
    WHERE c.champ_name = a.champ;
