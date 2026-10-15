package cargo

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/natanvictors/drafter/internal/parser"
)

const (
	apiURL = "https://lol.fandom.com/api.php"

	// season is the year embedded in the Leaguepedia page names we filter on
	// (e.g. "LCK/2026 Season/...").
	season = "2026"

	// gamesLimit is the max number of games returned per region. There is no
	// pagination yet, so only the most recent gamesLimit games are fetched.
	gamesLimit = 200

	maxRetries    = 3
	retryInterval = 5 // seconds
)

// numbered returns prefix1..prefixN (e.g. "Team1Pick1".."Team1Pick5").
func numbered(prefix string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = prefix + strconv.Itoa(i+1)
	}
	return out
}

// gameFields lists the Cargo fields requested for each game, in query order.
// They must stay in sync with the json tags of parser.Game.
func gameFields() []string {
	const pb = "PicksAndBansS7."

	fields := []string{
		"ScoreboardGames.Patch",
		pb + "_pageName=Page",
	}
	for _, prefix := range []string{
		"Team1Role", "Team2Role",
		"Team1Ban", "Team1Pick",
		"Team2Ban", "Team2Pick",
	} {
		for _, f := range numbered(prefix, 5) {
			fields = append(fields, pb+f)
		}
	}
	return append(fields,
		pb+"Team1",
		pb+"Team2",
		pb+"Winner",
		pb+"Team1PicksByRoleOrder",
		pb+"Team2PicksByRoleOrder",
		pb+"GameId",
		pb+"MatchId",
		"ScoreboardGames.DateTime_UTC",
	)
}

// gamesURL builds the Cargo query for the most recent finished games of a
// region (e.g. "LCK").
func gamesURL(region string) string {
	values := url.Values{
		"action":   {"cargoquery"},
		"format":   {"json"},
		"tables":   {"PicksAndBansS7,ScoreboardGames"},
		"fields":   {strings.Join(gameFields(), ",")},
		"join_on":  {"PicksAndBansS7.GameId=ScoreboardGames.GameId"},
		"order_by": {"ScoreboardGames.DateTime_UTC DESC"},
		"limit":    {strconv.Itoa(gamesLimit)},
	}
	values.Set("where", fmt.Sprintf(
		"PicksAndBansS7._pageName LIKE '%%%s/%s%%' AND ScoreboardGames.DateTime_UTC IS NOT NULL AND PicksAndBansS7.Winner != ''",
		region, season,
	))

	return apiURL + "?" + values.Encode()
}

// GamesByRegion fetches the picks, bans and results of the most recent games
// of a region (e.g. "LCK") from the Leaguepedia Cargo API.
func (c *Client) GamesByRegion(region string) (parser.Cargoresult, error) {
	return c.FetchAndParse(RequestInfo{
		URL:        gamesURL(region),
		Maxretries: maxRetries,
		Interval:   retryInterval,
	})
}
