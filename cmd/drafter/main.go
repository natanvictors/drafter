package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	cargo "github.com/natanvictors/drafter/internal/client"
	"github.com/natanvictors/drafter/internal/db"
	"github.com/natanvictors/drafter/internal/ingest"
)

var regions = []string{"CBLOL", "LCK", "LPL", "LEC", "LCS"}

func main() {
	godotenv.Load("../../.env")

	// TODO: skip regions updated recently, using UPDATE_TIME (seconds) and
	// queries.GetUpdateRecord(ctx, region) to read update_log.

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

	client, err := cargo.New(os.Getenv("API_BOT_NAME"), os.Getenv("API_BOT_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Login(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ingester := ingest.New(db.New(conn), client)

	if err := ingester.Champions(ctx); err != nil {
		log.Fatal(err)
	}

	for _, region := range regions {
		fmt.Println("Migrating games from", region)
		if err := ingester.Region(ctx, region); err != nil {
			log.Fatal(region, ": ", err)
		}
	}
}
