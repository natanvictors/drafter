package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/natanvictors/drafter/internal/db"
	"github.com/natanvictors/drafter/parser"
)

func main() {

	regions := []string{"CBLOL", "LCK", "LPL", "LEC", "LCS"}
	godotenv.Load("/home/natan/campify/.env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found on the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open connection to database")
	}
	defer conn.Close()

	queries := db.New(conn)

	ctx := context.Background()
	qtx := queries

	option := 2

	for option != 0 {
		fmt.Println("What would you like to do?")
		fmt.Println("1 - Load games into DB")
		fmt.Println("2 - Load champions into DB")
		fmt.Println("Press 0 to finish the program")

		fmt.Scan(&option)

		switch option {
		case 1:
			addGamesToDB(regions, ctx, *qtx)
		case 2:
			addChampionsToDB(ctx, *qtx)
		default:
			fmt.Println("Invalid option!")
		}

	}
}

func addGamesToDB(regions []string, ctx context.Context, qtx db.Queries) {
	for _, region := range regions {

		resp, err := os.ReadFile(fmt.Sprintf("../../regions_data/%s.json", region))
		if err != nil {
			log.Fatal("failed to read file:", err)
		}

		body, err := parser.ParseCargo(resp)
		if err != nil {
			log.Println("parse error:", err)
			continue
		}

		for _, game := range body.Cargoquery {
			_, err := qtx.CreateGame(ctx, db.CreateGameParams{
				Tournament:            game.Title.Page,
				Team1Role1:            game.Title.Team1Role1,
				Team1Role2:            game.Title.Team1Role2,
				Team1Role3:            game.Title.Team1Role3,
				Team1Role4:            game.Title.Team1Role4,
				Team1Role5:            game.Title.Team1Role5,
				Team2Role1:            game.Title.Team2Role1,
				Team2Role2:            game.Title.Team2Role2,
				Team2Role3:            game.Title.Team2Role3,
				Team2Role4:            game.Title.Team2Role4,
				Team2Role5:            game.Title.Team2Role5,
				Team1Ban1:             game.Title.Team1Ban1,
				Team1Ban2:             game.Title.Team1Ban2,
				Team1Ban3:             game.Title.Team1Ban3,
				Team1Ban4:             game.Title.Team1Ban4,
				Team1Ban5:             game.Title.Team1Ban5,
				Team1Pick1:            game.Title.Team1Pick1,
				Team1Pick2:            game.Title.Team1Pick2,
				Team1Pick3:            game.Title.Team1Pick3,
				Team1Pick4:            game.Title.Team1Pick4,
				Team1Pick5:            game.Title.Team1Pick5,
				Team2Ban1:             game.Title.Team2Ban1,
				Team2Ban2:             game.Title.Team2Ban2,
				Team2Ban3:             game.Title.Team2Ban3,
				Team2Ban4:             game.Title.Team2Ban4,
				Team2Ban5:             game.Title.Team2Ban5,
				Team2Pick1:            game.Title.Team2Pick1,
				Team2Pick2:            game.Title.Team2Pick2,
				Team2Pick3:            game.Title.Team2Pick3,
				Team2Pick4:            game.Title.Team2Pick4,
				Team2Pick5:            game.Title.Team2Pick5,
				Team1:                 game.Title.Team1,
				Team2:                 game.Title.Team2,
				Winner:                game.Title.Winner,
				Team1PicksByRoleOrder: game.Title.Team1PicksByRoleOrder,
				Team2PicksByRoleOrder: game.Title.Team2PicksByRoleOrder,
				GameID:                game.Title.GameID,
				MatchID:               game.Title.MatchID,
				DateTimeUtc:           game.Title.DateTimeUTC,
			})
			if err != nil {
				log.Fatal("insert failed: ", err)
			}
		}

		fmt.Println("Migrated succesfully!")
	}
}

func addChampionsToDB(ctx context.Context, qtx db.Queries) {
	resp, err := os.ReadFile("../../champions_data/champions.json")
	if err != nil {
		log.Fatal("failed to read file")
	}

	body, err := parser.ParseChampion(resp)
	if err != nil {
		log.Fatal("parse failed:", err)
	}

	fmt.Printf("Champions loaded: %d\n", len(body.Data))

	for _, raw := range body.Data {
		var name parser.Champion

		json.Unmarshal(raw, &name)

		_, err := qtx.AddChampions(ctx, name.Name)
		if err != nil {
			log.Fatal("insert failed: ", err)
		}

		fmt.Printf("Added champion %s\n", name)
	}

}
