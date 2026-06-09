package batch

import (
	"time"

	"github.com/natanvictors/drafter/internal/db"
)

type GameBatch struct {
	Tournament            []string
	Patch                 []string
	Team1Role1            []string
	Team1Role2            []string
	Team1Role3            []string
	Team1Role4            []string
	Team1Role5            []string
	Team2Role1            []string
	Team2Role2            []string
	Team2Role3            []string
	Team2Role4            []string
	Team2Role5            []string
	Team1Ban1             []string
	Team1Ban2             []string
	Team1Ban3             []string
	Team1Ban4             []string
	Team1Ban5             []string
	Team1Pick1            []string
	Team1Pick2            []string
	Team1Pick3            []string
	Team1Pick4            []string
	Team1Pick5            []string
	Team2Ban1             []string
	Team2Ban2             []string
	Team2Ban3             []string
	Team2Ban4             []string
	Team2Ban5             []string
	Team2Pick1            []string
	Team2Pick2            []string
	Team2Pick3            []string
	Team2Pick4            []string
	Team2Pick5            []string
	Team1                 []string
	Team2                 []string
	Winner                []string
	Team1PicksByRoleOrder []string
	Team2PicksByRoleOrder []string
	GameID                []string
	MatchID               []string
	DateTimeUTC           []time.Time
}

func (g GameBatch) ToSQLC() db.UpsertGamesParams {
	return db.UpsertGamesParams{
		Column1:  g.Tournament,
		Column2:  g.Patch,
		Column3:  g.Team1Role1,
		Column4:  g.Team1Role2,
		Column5:  g.Team1Role3,
		Column6:  g.Team1Role4,
		Column7:  g.Team1Role5,
		Column8:  g.Team2Role1,
		Column9:  g.Team2Role2,
		Column10: g.Team2Role3,
		Column11: g.Team2Role4,
		Column12: g.Team2Role5,
		Column13: g.Team1Ban1,
		Column14: g.Team1Ban2,
		Column15: g.Team1Ban3,
		Column16: g.Team1Ban4,
		Column17: g.Team1Ban5,
		Column18: g.Team1Pick1,
		Column19: g.Team1Pick2,
		Column20: g.Team1Pick3,
		Column21: g.Team1Pick4,
		Column22: g.Team1Pick5,
		Column23: g.Team2Ban1,
		Column24: g.Team2Ban2,
		Column25: g.Team2Ban3,
		Column26: g.Team2Ban4,
		Column27: g.Team2Ban5,
		Column28: g.Team2Pick1,
		Column29: g.Team2Pick2,
		Column30: g.Team2Pick3,
		Column31: g.Team2Pick4,
		Column32: g.Team2Pick5,
		Column33: g.Team1,
		Column34: g.Team2,
		Column35: g.Winner,
		Column36: g.Team1PicksByRoleOrder,
		Column37: g.Team2PicksByRoleOrder,
		Column38: g.GameID,
		Column39: g.MatchID,
		Column40: g.DateTimeUTC,
	}
}
