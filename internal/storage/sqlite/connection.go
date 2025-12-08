// Package sqlite implements a storage backend using SQLite.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver.
)

// NewConnection establishes a new SQLite database connection with the given path.
func NewConnection(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if pingErr := db.PingContext(ctx); pingErr != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil,
				fmt.Errorf("failed to ping database: %w, also failed to close db: %w", pingErr, closeErr)
		}
		return nil, pingErr
	}

	return db, nil
}
