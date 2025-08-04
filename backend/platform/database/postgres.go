package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/otterly-id/otterly/backend/pkg/utils"
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	postgresConnURL, err := utils.ConnectionURLBuilder("neon")
	if err != nil {
		return nil, fmt.Errorf("failed to build connection URL: %w", err)
	}

	db, err := sqlx.Connect("pgx", postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
