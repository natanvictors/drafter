-- name: UpsertGames :many
INSERT INTO games (
    tournament,
	patch,
    team1,
    team2,
    winner,
    game_id,
    match_id
)
SELECT
    unnest($1::text[]),
    unnest($2::text[]),
    unnest($3::int[]),
    unnest($4::int[]),
    unnest($5::text[]),
    unnest($6::text[]),
    unnest($7::text[])
ON CONFLICT (game_id) DO UPDATE SET
    tournament                = EXCLUDED.tournament,
	patch                     = EXCLUDED.patch,
    team1                     = EXCLUDED.team1,
    team2                     = EXCLUDED.team2,
    winner                    = EXCLUDED.winner,
    match_id                  = EXCLUDED.match_id
RETURNING ID;

-- name: UpsertGameTeams :many
INSERT INTO game_teams (
    side,
    name,
    roles,
    picks,
    bans
)
SELECT
    unnest($1::int[]),
    unnest($2::text[]),
    unnest($3::text[])::text[],
    unnest($4::text[])::text[],
    unnest($5::text[])::text[]
RETURNING ID;
