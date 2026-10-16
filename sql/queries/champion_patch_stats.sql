-- name: ClearChampionPatchStats :exec
DELETE FROM champion_patch_stats;

-- name: UpdateChampionPatchStats :exec
INSERT INTO champion_patch_stats (champ_id, role, patch, wins, losses, times_picked)
SELECT
    c.id,
    p.role,
    p.patch,
    COUNT(*) FILTER (WHERE p.won),
    COUNT(*) FILTER (WHERE NOT p.won),
    COUNT(*)
FROM (
    SELECT
        pk.champ,
        r.role,
        g.patch,
        gt.side::text = g.winner AS won
    FROM games g
    JOIN game_teams gt ON gt.id IN (g.team1, g.team2)
    CROSS JOIN LATERAL unnest(gt.picks) WITH ORDINALITY AS pk(champ, idx)
    JOIN LATERAL unnest(gt.roles) WITH ORDINALITY AS r(role, idx) ON r.idx = pk.idx
    WHERE g.winner IN ('1', '2')
      AND g.patch ~ '^[0-9]+\.[0-9]+$'
) p
JOIN champions c ON c.champ_name = p.champ
GROUP BY c.id, p.role, p.patch;