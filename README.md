# Drafter

This is a LoL drafter API.

## Project layout

```
cmd/drafter        entrypoint: config, DB connection, login, runs the ingestion
internal/ingest    orchestration: fetch -> map -> upsert (Champions, Region)
internal/client    Leaguepedia/Cargo HTTP client (login, retries, GamesByRegion)
internal/ddragon   Data Dragon client (champion names)
internal/parser    JSON structs for Cargo / Data Dragon responses
internal/db        sqlc-generated code (do not edit) + goose migrations runner
internal/db/batch  column-oriented batches for the unnest-based bulk upserts
sql/schema         goose migrations
sql/queries        sqlc queries
```

Data flow:

```
Leaguepedia (Cargo) -> client.GamesByRegion -> parser.Game -> ingest.Region
                                                  -> batch.GameTeamBatch / GameBatch
                                                  -> db (game_teams, games, champions_data)
Data Dragon -> ddragon.ChampionNames -> ingest.Champions -> db (champions)
```
