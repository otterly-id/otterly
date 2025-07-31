package database

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/otterly-id/otterly/backend/pkg/utils"
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	maxConn := utils.GetEnvAsInt("DB_MAX_CONNECTIONS", 25)
	maxIdleConn := utils.GetEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5)
	maxLifetimeConn := utils.GetEnvAsInt("DB_MAX_LIFETIME_CONNECTIONS", 300)

	postgresConnURL, err := utils.ConnectionURLBuilder("neon")
	if err != nil {
		return nil, fmt.Errorf("failed to build connection URL: %w", err)
	}

	db, err := sqlx.Connect("pgx", postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn) * time.Second)

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
