-- +goose Up

CREATE TABLE games (
	ID SERIAL NOT NULL PRIMARY KEY,
	patch TEXT,
	tournament TEXT NOT NULL,
	team1 TEXT NOT NULL,
	team2 TEXT NOT NULL,
	winner TEXT NOT NULL,
	game_id TEXT UNIQUE NOT NULL,
	match_id TEXT NOT NULL,
	date_time_utc TIMESTAMPTZ NOT NULL
);

CREATE TABLE game_teams (
	ID SERIAL NOT NULL PRIMARY KEY, 
	game_id INTEGER NOT NULL REFERENCES games(ID),
	side INT NOT NULL,
	name TEXT NOT NULL,
	roles TEXT[5],
	picks TEXT[5],
	bans TEXT[5],
	UNIQUE (game_id, side)
);


-- +goose Down
DROP TABLE game_teams;
DROP TABLE games;
