-- name: ClearChampionStats :exec
DELETE FROM champions_data;

-- name: UpdateChampionStats :exec

WITH picks AS (
    SELECT champ, role, won FROM (
        SELECT team1_pick1 AS champ, team1_role1 AS role, (winner='1') AS won FROM games
        UNION ALL
        SELECT team1_pick2, team1_role2, (winner = '1') FROM games
        UNION ALL
        SELECT team1_pick3, team1_role3, (winner = '1') FROM games
        UNION ALL
        SELECT team1_pick4, team1_role4, (winner = '1') FROM games
        UNION ALL
        SELECT team1_pick5, team1_role5, (winner = '1') FROM games
        UNION ALL
        SELECT team2_pick1, team2_role1, (winner = '2') FROM games
        UNION ALL
        SELECT team2_pick2, team2_role2, (winner = '2') FROM games
        UNION ALL
        SELECT team2_pick3, team2_role3, (winner = '2') FROM games
        UNION ALL
        SELECT team2_pick4, team2_role4, (winner = '2') FROM games
        UNION ALL
        SELECT team2_pick5, team2_role5, (winner = '2') FROM games
    )sub
),
agg AS (
    SELECT champ,
    COUNT(*) AS times_picked,
    COUNT(*) FILTER (WHERE role = 'Top') AS picked_top,
    COUNT(*) FILTER (WHERE role = 'Mid') AS picked_mid,
    COUNT(*) FILTER (WHERE role = 'Bot') AS picked_adc,
    COUNT(*) FILTER (WHERE role = 'Support') AS picked_sup,
    COUNT(*) FILTER (WHERE role = 'Jungle') AS picked_jungle,
    COUNT(*) FILTER (WHERE role = 'Top' AND won) AS wins_top,
    COUNT(*) FILTER (WHERE role = 'Mid' AND won) AS wins_mid,
    COUNT(*) FILTER (WHERE role = 'Bot' AND won) AS wins_adc,
    COUNT(*) FILTER (WHERE role = 'Support' AND won) AS wins_sup,
    COUNT(*) FILTER (WHERE role = 'Jungle' AND won) AS wins_jg,
    COUNT(*) FILTER (WHERE role = 'Top' AND NOT won) AS losses_top,
    COUNT(*) FILTER (WHERE role = 'Mid' AND NOT won) AS losses_mid,
    COUNT(*) FILTER (WHERE role = 'Bot' AND NOT won) AS losses_adc,
    COUNT(*) FILTER (WHERE role = 'Support' AND NOT won) AS losses_sup,
    COUNT(*) FILTER (WHERE role = 'Jungle' AND NOT won) AS losses_jg
    FROM picks
    GROUP BY champ
)INSERT INTO champions_data(
    champ_id,
    times_picked,
    picked_top,
    picked_mid,
    picked_adc,
    picked_sup,
    picked_jungle,
    wins_top,
    wins_mid,
    wins_jg,
    wins_adc,
    wins_sup,
    losses_top,
    losses_jg,
    losses_mid,
    losses_adc,
    losses_sup
)SELECT
    c.id,
    a.times_picked,
    a.picked_top,
    a.picked_mid,
    a.picked_adc,
    a.picked_sup,
    a.picked_jungle,
    a.wins_top,
    a.wins_mid,
    a.wins_jg,
    a.wins_adc,
    a.wins_sup,
    a.losses_top,
    a.losses_jg,
    a.losses_mid,
    a.losses_adc,
    a.losses_sup
FROM agg a
JOIN champions c ON c.champ_name = a.champ;
