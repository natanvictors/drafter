-- +goose Up

CREATE TABLE game_teams (
	ID SERIAL NOT NULL PRIMARY KEY, 
	side INTEGER NOT NULL,
	name TEXT NOT NULL,
	roles TEXT[5],
	picks TEXT[5],
	bans TEXT[5]
);

CREATE TABLE games (
	ID SERIAL NOT NULL PRIMARY KEY,
	patch TEXT,
	tournament TEXT,
	team1 INTEGER REFERENCES game_teams(ID),
	team2 INTEGER REFERENCES game_teams(ID),
	winner TEXT NOT NULL,
	game_id TEXT UNIQUE NOT NULL,
	match_id TEXT NOT NULL
);

-- +goose Down
DROP TABLE games;
DROP TABLE game_teams;
