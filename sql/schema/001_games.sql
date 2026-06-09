-- +goose Up

CREATE TABLE games (
	ID SERIAL NOT NULL PRIMARY KEY,
	patch TEXT,
	tournament TEXT NOT NULL,
	team1_role1 TEXT NOT NULL,
	team1_role2 TEXT NOT NULL,
	team1_role3 TEXT NOT NULL,
	team1_role4 TEXT NOT NULL,
	team1_role5 TEXT NOT NULL,
	team2_role1 TEXT NOT NULL,
	team2_role2 TEXT NOT NULL,
	team2_role3 TEXT NOT NULL,
	team2_role4 TEXT NOT NULL,
	team2_role5 TEXT NOT NULL,
	team1_ban1 TEXT NOT NULL,
	team1_ban2 TEXT NOT NULL,
	team1_ban3 TEXT NOT NULL,
	team1_ban4 TEXT NOT NULL,
	team1_ban5 TEXT NOT NULL,
	team1_pick1 TEXT NOT NULL,
	team1_pick2 TEXT NOT NULL,
	team1_pick3 TEXT NOT NULL,
	team1_pick4 TEXT NOT NULL,
	team1_pick5 TEXT NOT NULL,
	team2_ban1 TEXT NOT NULL,
	team2_ban2 TEXT NOT NULL,
	team2_ban3 TEXT NOT NULL,
	team2_ban4 TEXT NOT NULL,
	team2_ban5 TEXT NOT NULL,
	team2_pick1 TEXT NOT NULL,
	team2_pick2 TEXT NOT NULL,
	team2_pick3 TEXT NOT NULL,
	team2_pick4 TEXT NOT NULL,
	team2_pick5 TEXT NOT NULL,
	team1 TEXT NOT NULL,
	team2 TEXT NOT NULL,
	winner TEXT NOT NULL,
	team1_picks_by_role_order TEXT NOT NULL,
	team2_picks_by_role_order TEXT NOT NULL,
	game_id TEXT UNIQUE NOT NULL,
	match_id TEXT NOT NULL,
	date_time_utc TIMESTAMPTZ NOT NULL
);


-- +goose Down

DROP TABLE games;