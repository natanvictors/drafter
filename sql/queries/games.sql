-- name: UpsertGames :many
INSERT INTO games (
    tournament,
	patch,
    team1,
    team2,
    winner,
    game_id,
    match_id,
    date_time_utc
)
SELECT
    unnest($1::text[]),
    unnest($2::text[]),
    unnest($3::text[]),
    unnest($4::text[]),
    unnest($5::text[]),
    unnest($6::text[]),
    unnest($7::text[]),
    unnest($8::timestamptz[])
ON CONFLICT (game_id) DO UPDATE SET
    tournament                = EXCLUDED.tournament,
	patch                     = EXCLUDED.patch,
    team1                     = EXCLUDED.team1,
    team2                     = EXCLUDED.team2,
    winner                    = EXCLUDED.winner,
    match_id                  = EXCLUDED.match_id,
    date_time_utc             = EXCLUDED.date_time_utc
RETURNING ID, game_id;

-- name: UpsertGameTeams :exec
INSERT INTO game_teams (
    game_id,
    side,
    name,
    roles,
    picks,
    bans
)
SELECT
    unnest($1::int[]),
    unnest($2::int[]),
    unnest($3::text[]),
    unnest($4::text[][]),
    unnest($5::text[][]),
    unnest($6::text[][])
ON CONFLICT (game_id, side) DO UPDATE SET
    name = EXCLUDED.name,
    roles = EXCLUDED.roles,
    picks = EXCLUDED.picks,
    bans = EXCLUDED.bans;