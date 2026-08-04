package batch

import (
	"github.com/natanvictors/drafter/internal/db"
)

type GameBatch struct {
	Tournament []string
	Patch      []string
	Team1      []int32
	Team2      []int32
	Winner     []string
	GameID     []string
	MatchID    []string
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
	}
}

type GameTeamBatch struct {
	Side []int32
	Name []string
	// Roles, Picks, and Bans must be Postgres array literals (e.g. "{Top,Jungle,...}"),
	// not raw values — UpsertGameTeams unnests one literal per team and casts it back to text[].
	Roles []string
	Picks []string
	Bans  []string
}

func (t GameTeamBatch) ToSQLC() db.UpsertGameTeamsParams {
	return db.UpsertGameTeamsParams{
		Column1: t.Side,
		Column2: t.Name,
		Column3: t.Roles,
		Column4: t.Picks,
		Column5: t.Bans,
	}
}
