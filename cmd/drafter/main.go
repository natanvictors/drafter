package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	cargo "github.com/natanvictors/drafter/internal/client"
	"github.com/natanvictors/drafter/internal/db"
	"github.com/natanvictors/drafter/internal/parser"
)

func main() {

	regions := []string{"CBLOL", "LCK", "LPL", "LEC", "LCS"}
	godotenv.Load("../../.env")

	// updateTime, err := strconv.Atoi(os.Getenv("UPDATE_TIME"))
	// if err != nil {
	// 	log.Fatal("Failed to get envvar UPDATE_TIME")
	// }

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found on the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open connection to database")
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	queries := db.New(conn)

	ctx := context.Background()
	qtx := queries

	client, err := cargo.New(os.Getenv("API_BOT_NAME"), os.Getenv("API_BOT_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Login()
	if err != nil {
		log.Fatal(err)
	}
	//updateInterval := time.Duration(updateTime) * time.Second

	addChampionsToDB(ctx, *qtx)
	for _, region := range regions {
		fmt.Println("Migrating games from", region)
		//record, err := qtx.GetUpdateRecord(ctx, region)
		// timeUntilUpdate := updateInterval - time.Since(record.LastUpdatedAt)
		// if err == nil && timeUntilUpdate > 0 {
		// 	fmt.Printf("Skipping %s, updated recently\nTime until next update: %x\n", region, timeUntilUpdate)
		// 	continue
		// }

		updateChampionData(ctx, qtx, region, client)
	}

}

// arrayLiteral encodes a Go string slice as a Postgres array literal (e.g. `{a,b,c}`)
// so it can round-trip through a text[] bulk-insert param and be cast back with ::text[].
func arrayLiteral(ss []string) string {
	val, err := pq.StringArray(ss).Value()
	if err != nil {
		log.Fatal("failed to encode array literal: ", err)
	}
	return val.(string)
}

func updateChampionData(ctx context.Context, qtx *db.Queries, region string, c *cargo.Client) {

	urlValues := url.Values{
		"action":   {"cargoquery"},
		"format":   {"json"},
		"tables":   {"PicksAndBansS7,ScoreboardGames"},
		"fields":   {"ScoreboardGames.Patch,PicksAndBansS7._pageName=Page,PicksAndBansS7.Team1Role1,PicksAndBansS7.Team1Role2,PicksAndBansS7.Team1Role3,PicksAndBansS7.Team1Role4,PicksAndBansS7.Team1Role5,PicksAndBansS7.Team2Role1,PicksAndBansS7.Team2Role2,PicksAndBansS7.Team2Role3,PicksAndBansS7.Team2Role4,PicksAndBansS7.Team2Role5,PicksAndBansS7.Team1Ban1,PicksAndBansS7.Team1Ban2,PicksAndBansS7.Team1Ban3,PicksAndBansS7.Team1Ban4,PicksAndBansS7.Team1Ban5,PicksAndBansS7.Team1Pick1,PicksAndBansS7.Team1Pick2,PicksAndBansS7.Team1Pick3,PicksAndBansS7.Team1Pick4,PicksAndBansS7.Team1Pick5,PicksAndBansS7.Team2Ban1,PicksAndBansS7.Team2Ban2,PicksAndBansS7.Team2Ban3,PicksAndBansS7.Team2Ban4,PicksAndBansS7.Team2Ban5,PicksAndBansS7.Team2Pick1,PicksAndBansS7.Team2Pick2,PicksAndBansS7.Team2Pick3,PicksAndBansS7.Team2Pick4,PicksAndBansS7.Team2Pick5,PicksAndBansS7.Team1,PicksAndBansS7.Team2,PicksAndBansS7.Winner,PicksAndBansS7.Team1PicksByRoleOrder,PicksAndBansS7.Team2PicksByRoleOrder,PicksAndBansS7.GameId,PicksAndBansS7.MatchId,ScoreboardGames.DateTime_UTC"},
		"join_on":  {"PicksAndBansS7.GameId=ScoreboardGames.GameId"},
		"order_by": {"ScoreboardGames.DateTime_UTC DESC"},
		"limit":    {"200"},
	}

	urlValues.Set("where", fmt.Sprintf("PicksAndBansS7._pageName LIKE '%%%s/2026%%' AND ScoreboardGames.DateTime_UTC IS NOT NULL AND PicksAndBansS7.Winner != ''", region))

	info := cargo.RequestInfo{
		URL:        "https://lol.fandom.com/api.php?" + urlValues.Encode(),
		Maxretries: 3,
		Interval:   5,
	}
	//log.Print(info.URL)
	body, err := c.FetchAndParse(info)
	if err != nil {
		log.Fatal("fetch failed: ", err)
	}

	log.Printf("Games found for %s: %d\n", region, len(body.Cargoresponse))

	fmt.Printf("Migrated %s succesfully!\n", region)

	gamesParams := db.UpsertGamesParams{}

	teamsParams := db.UpsertGameTeamsParams{}

	for _, game := range body.Cargoresponse {
		t := game.Game
		team1, team2 := t.Teams()
		if t.Championship == "" {
			log.Printf("empty: %v", t)
		}
		teamsParams.Column1 = append(teamsParams.Column1, int32(team1.Side), int32(team2.Side))
		teamsParams.Column2 = append(teamsParams.Column2, team1.Name, team2.Name)
		teamsParams.Column3 = append(teamsParams.Column3, arrayLiteral(team1.Roles), arrayLiteral(team2.Roles))
		teamsParams.Column4 = append(teamsParams.Column4, arrayLiteral(team1.Picks), arrayLiteral(team2.Picks))
		teamsParams.Column5 = append(teamsParams.Column5, arrayLiteral(team1.Bans), arrayLiteral(team2.Bans))

		gamesParams.Column1 = append(gamesParams.Column1, t.Championship)
		gamesParams.Column2 = append(gamesParams.Column2, t.Patch)
		gamesParams.Column5 = append(gamesParams.Column5, t.Winner)
		gamesParams.Column6 = append(gamesParams.Column6, t.GameID)
		gamesParams.Column7 = append(gamesParams.Column7, t.MatchID)
	}

	teamIds, err := qtx.UpsertGameTeams(ctx, teamsParams)
	if err != nil {
		log.Fatal("UpsertTeamsParams failed: ", err)
	}

	team1Ids := make([]int32, 0, len(teamIds)/2)
	team2Ids := make([]int32, 0, len(teamIds)/2)

	for idx, el := range teamIds {
		if idx%2 == 0 {
			team1Ids = append(team1Ids, el)
		} else {
			team2Ids = append(team2Ids, el)
		}
	}
	gamesParams.Column3 = append(gamesParams.Column3, team1Ids...)
	gamesParams.Column4 = append(gamesParams.Column4, team2Ids...)

	_, err = qtx.UpsertGames(ctx, gamesParams)
	if err != nil {
		log.Println("UpsertGames failed:", err)
	}

	err = qtx.UpsertRecord(ctx, db.UpsertRecordParams{
		Region:        region,
		LastUpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("failed to update record:", err)
	}
	addWinsAndLoses(ctx, *qtx)
}

func addChampionsToDB(ctx context.Context, qtx db.Queries) {
	resp, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		log.Fatal("failed to get version")
	}
	defer resp.Body.Close()

	var versions []string

	err = json.NewDecoder(resp.Body).Decode(&versions)
	if err != nil {
		log.Fatal("failed to decode versions API")
	}

	resp.Body.Close()

	resp, err = http.Get(fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json", versions[0]))
	if err != nil {
		log.Fatal("failed to get champion API")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read failed:", err)
	}

	body, err := parser.ParseChampion(data)
	if err != nil {
		log.Fatal("failed to parse data:", err)
	}

	fmt.Printf("Champions loaded: %d\n", len(body.Data))

	for _, raw := range body.Data {
		var name parser.Champion

		json.Unmarshal(raw, &name)

		_, err := qtx.AddChampions(ctx, name.Name)
		if err != nil {
			log.Fatal("insert failed: ", err)
		}
	}

}

func addWinsAndLoses(ctx context.Context, qtx db.Queries) {

	err := qtx.ClearChampionStats(ctx)
	if err != nil {
		log.Fatal("failed to clear champion stats:", err)
	}

	if err := qtx.UpdateChampionStats(ctx); err != nil {
		log.Fatal("failed to update champion stats:", err)
	}

	fmt.Println("Updated succesfully!")
}
