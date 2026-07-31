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
	_ "github.com/lib/pq"
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

func updateChampionData(ctx context.Context, qtx *db.Queries, region string, c *cargo.Client) {

	urlValues := url.Values{
		"action":   {"cargoquery"},
		"format":   {"json"},
		"tables":   {"PicksAndBansS7,ScoreboardGames"},
		"fields":   {"ScoreboardGames.Patch,PicksAndBansS7._pageName=Page,PicksAndBansS7.Team1Role1,PicksAndBansS7.Team1Role2,PicksAndBansS7.Team1Role3,PicksAndBansS7.Team1Role4,PicksAndBansS7.Team1Role5,PicksAndBansS7.Team2Role1,PicksAndBansS7.Team2Role2,PicksAndBansS7.Team2Role3,PicksAndBansS7.Team2Role4,PicksAndBansS7.Team2Role5,PicksAndBansS7.Team1Ban1,PicksAndBansS7.Team1Ban2,PicksAndBansS7.Team1Ban3,PicksAndBansS7.Team1Ban4,PicksAndBansS7.Team1Ban5,PicksAndBansS7.Team1Pick1,PicksAndBansS7.Team1Pick2,PicksAndBansS7.Team1Pick3,PicksAndBansS7.Team1Pick4,PicksAndBansS7.Team1Pick5,PicksAndBansS7.Team2Ban1,PicksAndBansS7.Team2Ban2,PicksAndBansS7.Team2Ban3,PicksAndBansS7.Team2Ban4,PicksAndBansS7.Team2Ban5,PicksAndBansS7.Team2Pick1,PicksAndBansS7.Team2Pick2,PicksAndBansS7.Team2Pick3,PicksAndBansS7.Team2Pick4,PicksAndBansS7.Team2Pick5,PicksAndBansS7.Team1,PicksAndBansS7.Team2,PicksAndBansS7.Winner,PicksAndBansS7.Team1PicksByRoleOrder,PicksAndBansS7.Team2PicksByRoleOrder,PicksAndBansS7.GameId,PicksAndBansS7.MatchId,ScoreboardGames.DateTime_UTC"},
		"join_on":  {"PicksAndBansS7.GameId=ScoreboardGames.GameId"},
		"order_by": {"ScoreboardGames.DateTime_UTC DESC"},
		"limit":    {"1000"},
	}

	urlValues.Set("where", fmt.Sprintf("PicksAndBansS7._pageName LIKE '%%%s/2026%%' AND ScoreboardGames.DateTime_UTC IS NOT NULL AND PicksAndBansS7.Winner != ''", region))

	info := cargo.RequestInfo{
		URL:        "https://lol.fandom.com/api.php?" + urlValues.Encode(),
		Maxretries: 3,
		Interval:   5,
	}

	body, err := c.FetchAndParse(info)
	if err != nil {
		log.Fatal("fetch failed: ", err)
	}

	log.Printf("Games found for %s: %d\n", region, len(body.Cargoresponse))

	fmt.Printf("Migrated %s succesfully!\n", region)

	// place upsertgames and upsertteamgames here

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
