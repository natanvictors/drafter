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
	"github.com/natanvictors/drafter/internal/client"
	"github.com/natanvictors/drafter/internal/db"
	"github.com/natanvictors/drafter/internal/parser"
)

func main() {

	regions := []string{"CBLOL", "LCK", "LPL", "LEC", "LCS"}
	godotenv.Load("../../.env")

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

	client, err := client.New()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Login()
	if err != nil {
		log.Fatal(err)
	}

	addChampionsToDB(ctx, *qtx)

	for _, region := range regions {
		record, err := qtx.GetUpdateRecord(ctx, region)
		if err == nil && time.Since(record.LastUpdatedAt) < 12*time.Hour {
			fmt.Printf("Skipping %s, updated recently\n", region)
			continue
		}

		updateChampionData(ctx, qtx, region, client)
	}

}

func updateChampionData(ctx context.Context, qtx *db.Queries, region string, c *client.Client) {

	urlValues := url.Values{
		"action":   {"cargoquery"},
		"format":   {"json"},
		"tables":   {"PicksAndBansS7,ScoreboardGames"},
		"fields":   {"ScoreboardGames.Patch,PicksAndBansS7._pageName=Page,PicksAndBansS7.Team1Role1,PicksAndBansS7.Team1Role2,PicksAndBansS7.Team1Role3,PicksAndBansS7.Team1Role4,PicksAndBansS7.Team1Role5,PicksAndBansS7.Team2Role1,PicksAndBansS7.Team2Role2,PicksAndBansS7.Team2Role3,PicksAndBansS7.Team2Role4,PicksAndBansS7.Team2Role5,PicksAndBansS7.Team1Ban1,PicksAndBansS7.Team1Ban2,PicksAndBansS7.Team1Ban3,PicksAndBansS7.Team1Ban4,PicksAndBansS7.Team1Ban5,PicksAndBansS7.Team1Pick1,PicksAndBansS7.Team1Pick2,PicksAndBansS7.Team1Pick3,PicksAndBansS7.Team1Pick4,PicksAndBansS7.Team1Pick5,PicksAndBansS7.Team2Ban1,PicksAndBansS7.Team2Ban2,PicksAndBansS7.Team2Ban3,PicksAndBansS7.Team2Ban4,PicksAndBansS7.Team2Ban5,PicksAndBansS7.Team2Pick1,PicksAndBansS7.Team2Pick2,PicksAndBansS7.Team2Pick3,PicksAndBansS7.Team2Pick4,PicksAndBansS7.Team2Pick5,PicksAndBansS7.Team1,PicksAndBansS7.Team2,PicksAndBansS7.Winner,PicksAndBansS7.Team1PicksByRoleOrder,PicksAndBansS7.Team2PicksByRoleOrder,PicksAndBansS7.GameId,PicksAndBansS7.MatchId,ScoreboardGames.DateTime_UTC"},
		"join_on":  {"PicksAndBansS7.GameId=ScoreboardGames.GameId"},
		"order_by": {"ScoreboardGames.DateTime_UTC DESC"},
		"limit":    {"500"},
	}

	urlValues.Set("where", fmt.Sprintf("PicksAndBansS7._pageName LIKE '%%%s/2026%%' AND ScoreboardGames.DateTime_UTC IS NOT NULL AND PicksAndBansS7.Winner != ''", region))
	resp, err := c.Fetch("https://lol.fandom.com/api.php?" + urlValues.Encode())
	if err != nil {
		log.Println(err)
	}
	body, err := parser.ParseCargo(resp)
	if err != nil {
		log.Println("parse error:", err)
	}

	log.Printf("Games found for %s: %d\n", region, len(body.Cargoquery))

	params := db.UpsertGamesParams{}

	for _, game := range body.Cargoquery {
		t := game.Title

		log.Printf("Page: %q", t.Page)

		params.Column1 = append(params.Column1, t.Page)
		params.Column2 = append(params.Column2, t.Patch)
		params.Column3 = append(params.Column3, t.Team1Role1)
		params.Column4 = append(params.Column4, t.Team1Role2)
		params.Column5 = append(params.Column5, t.Team1Role3)
		params.Column6 = append(params.Column6, t.Team1Role4)
		params.Column7 = append(params.Column7, t.Team1Role5)
		params.Column8 = append(params.Column8, t.Team2Role1)
		params.Column9 = append(params.Column9, t.Team2Role2)
		params.Column10 = append(params.Column10, t.Team2Role3)
		params.Column11 = append(params.Column11, t.Team2Role4)
		params.Column12 = append(params.Column12, t.Team2Role5)
		params.Column13 = append(params.Column13, t.Team1Ban1)
		params.Column14 = append(params.Column14, t.Team1Ban2)
		params.Column15 = append(params.Column15, t.Team1Ban3)
		params.Column16 = append(params.Column16, t.Team1Ban4)
		params.Column17 = append(params.Column17, t.Team1Ban5)
		params.Column18 = append(params.Column18, t.Team1Pick1)
		params.Column19 = append(params.Column19, t.Team1Pick2)
		params.Column20 = append(params.Column20, t.Team1Pick3)
		params.Column21 = append(params.Column21, t.Team1Pick4)
		params.Column22 = append(params.Column22, t.Team1Pick5)
		params.Column23 = append(params.Column23, t.Team2Ban1)
		params.Column24 = append(params.Column24, t.Team2Ban2)
		params.Column25 = append(params.Column25, t.Team2Ban3)
		params.Column26 = append(params.Column26, t.Team2Ban4)
		params.Column27 = append(params.Column27, t.Team2Ban5)
		params.Column28 = append(params.Column28, t.Team2Pick1)
		params.Column29 = append(params.Column29, t.Team2Pick2)
		params.Column30 = append(params.Column30, t.Team2Pick3)
		params.Column31 = append(params.Column31, t.Team2Pick4)
		params.Column32 = append(params.Column32, t.Team2Pick5)
		params.Column33 = append(params.Column33, t.Team1)
		params.Column34 = append(params.Column34, t.Team2)
		params.Column35 = append(params.Column35, t.Winner)
		params.Column36 = append(params.Column36, t.Team1PicksByRoleOrder)
		params.Column37 = append(params.Column37, t.Team2PicksByRoleOrder)
		params.Column38 = append(params.Column38, t.GameID)
		params.Column39 = append(params.Column39, t.MatchID)

		parsedTime, err := time.Parse("2006-01-02 15:04:05", t.DateTimeUTC)
		if err != nil {
			log.Println("failed to parse time", err)
		}

		params.Column40 = append(params.Column40, parsedTime)
	}

	err = qtx.UpsertGames(ctx, params)
	if err != nil {
		log.Println("upsert failed: ", err)
	}

	fmt.Printf("Migrated %s succesfully!\n", region)

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
