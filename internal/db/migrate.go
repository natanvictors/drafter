package db

import (
	"database/sql"
	"fmt"

	"github.com/natanvictors/drafter/sql/schema"
	"github.com/pressly/goose/v3"
)

// Migrate applies any pending schema migrations, creating tables that don't
// exist yet. It's safe to call on every startup.
func Migrate(conn *sql.DB) error {
	goose.SetBaseFS(schema.FS)
	defer goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(conn, "."); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
