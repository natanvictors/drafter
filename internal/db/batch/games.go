package batch

import (
	"time"

	"github.com/natanvictors/drafter/internal/db"
)

type GameBatch struct {
	Tournament  []string
	Patch       []string
	Team1       []string
	Team2       []string
	Winner      []string
	GameID      []string
	MatchID     []string
	DateTimeUTC []time.Time
}

func (g GameBatch) ToSQLC() db.UpsertGamesParams {
	return db.UpsertGamesParams{
		Column1: g.Tournament,
		Column2: g.Patch,
		Column3: g.Team1,
		Column4: g.Team2,
		Column5: g.Winner,
		Column6: g.GameID,
		Column7: g.MatchID,
		Column8: g.DateTimeUTC,
	}
}

type GameTeamBatch struct {
	GameID []int32
	Side   []int32
	Name   []string
	Roles  [][]string
	Picks  [][]string
	Bans   [][]string
}

func (t GameTeamBatch) ToSQLC() db.UpsertGameTeamsParams {
	return db.UpsertGameTeamsParams{
		Column1: t.GameID,
		Column2: t.Side,
		Column3: t.Name,
		Column4: t.Roles,
		Column5: t.Picks,
		Column6: t.Bans,
	}
}
