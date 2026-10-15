// Package ingest moves data from the external sources (Leaguepedia, Data
// Dragon) into the database.
package ingest

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	cargo "github.com/natanvictors/drafter/internal/client"
	"github.com/natanvictors/drafter/internal/db"
	"github.com/natanvictors/drafter/internal/db/batch"
	"github.com/natanvictors/drafter/internal/ddragon"
)

// Ingester orchestrates fetching external data and persisting it.
type Ingester struct {
	queries *db.Queries
	cargo   *cargo.Client
}

func New(queries *db.Queries, c *cargo.Client) *Ingester {
	return &Ingester{queries: queries, cargo: c}
}

// Champions loads every champion name from Data Dragon into the database.
func (i *Ingester) Champions(ctx context.Context) error {
	names, err := ddragon.ChampionNames()
	if err != nil {
		return err
	}

	fmt.Printf("Champions loaded: %d\n", len(names))

	for _, name := range names {
		if _, err := i.queries.AddChampions(ctx, name); err != nil {
			return fmt.Errorf("insert failed: %w", err)
		}
	}
	return nil
}

// Region fetches the latest games of a region, upserts them (and their
// per-side draft data) and refreshes the champion stats.
func (i *Ingester) Region(ctx context.Context, region string) error {
	body, err := i.cargo.GamesByRegion(region)
	if err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}

	log.Printf("Games found for %s: %d\n", region, len(body.Cargoresponse))
	fmt.Printf("Migrated %s succesfully!\n", region)

	var games batch.GameBatch
	var teams batch.GameTeamBatch

	for _, game := range body.Cargoresponse {
		g := game.Game
		if g.Championship == "" {
			log.Printf("empty: %v", g)
		}

		team1, team2 := g.Teams()
		for _, t := range []struct {
			name               string
			side               int
			roles, picks, bans []string
		}{
			{team1.Name, team1.Side, team1.Roles, team1.Picks, team1.Bans},
			{team2.Name, team2.Side, team2.Roles, team2.Picks, team2.Bans},
		} {
			teams.Side = append(teams.Side, int32(t.side))
			teams.Name = append(teams.Name, t.name)
			teams.Roles = append(teams.Roles, arrayLiteral(t.roles))
			teams.Picks = append(teams.Picks, arrayLiteral(t.picks))
			teams.Bans = append(teams.Bans, arrayLiteral(t.bans))
		}

		games.Tournament = append(games.Tournament, g.Championship)
		games.Patch = append(games.Patch, g.Patch)
		games.Winner = append(games.Winner, g.Winner)
		games.GameID = append(games.GameID, g.GameID)
		games.MatchID = append(games.MatchID, g.MatchID)
	}

	// UpsertGameTeams returns two IDs per game (team1, team2), in insertion order.
	teamIDs, err := i.queries.UpsertGameTeams(ctx, teams.ToSQLC())
	if err != nil {
		return fmt.Errorf("UpsertTeamsParams failed: %w", err)
	}
	for idx, id := range teamIDs {
		if idx%2 == 0 {
			games.Team1 = append(games.Team1, id)
		} else {
			games.Team2 = append(games.Team2, id)
		}
	}

	if _, err := i.queries.UpsertGames(ctx, games.ToSQLC()); err != nil {
		log.Println("UpsertGames failed:", err)
	}

	if err := i.queries.DeleteOrphanedGameTeams(ctx); err != nil {
		log.Println("failed to clean up orphaned game_teams:", err)
	}

	err = i.queries.UpsertRecord(ctx, db.UpsertRecordParams{
		Region:        region,
		LastUpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("failed to update record:", err)
	}

	return i.refreshChampionStats(ctx)
}

// refreshChampionStats recomputes the per-champion pick/win/loss counters.
func (i *Ingester) refreshChampionStats(ctx context.Context) error {
	if err := i.queries.ClearChampionStats(ctx); err != nil {
		return fmt.Errorf("failed to clear champion stats: %w", err)
	}
	if err := i.queries.UpdateChampionStats(ctx); err != nil {
		return fmt.Errorf("failed to update champion stats: %w", err)
	}

	fmt.Println("Updated succesfully!")
	return nil
}

// arrayLiteral encodes a Go string slice as a Postgres array literal (e.g. `{a,b,c}`)
// so it can round-trip through a text[] bulk-insert param and be cast back with ::text[].
func arrayLiteral(ss []string) string {
	val, err := pq.StringArray(ss).Value()
	if err != nil {
		// pq.StringArray.Value never fails for []string; keep the guard visible.
		panic(fmt.Sprintf("failed to encode array literal: %v", err))
	}
	return val.(string)
}
