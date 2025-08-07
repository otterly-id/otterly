package db

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/otterly-id/otterly/backend/internal/api/queries"
	"github.com/spf13/viper"
)

var (
	dbInstance *Queries
	mu         sync.RWMutex
)

type Queries struct {
	*queries.UserQueries
	*queries.AuthQueries
}

func PostgreSQLConnection(config *viper.Viper) (*sqlx.DB, error) {
	postgresConnURL := config.GetString("DB_URL")

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

func GetDBConnection(config *viper.Viper) (*Queries, error) {
	mu.RLock()
	if dbInstance != nil {
		mu.RUnlock()
		return dbInstance, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if dbInstance != nil {
		return dbInstance, nil
	}

	db, err := PostgreSQLConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	setupConnectionPool(db, config)

	dbInstance = &Queries{
		UserQueries: &queries.UserQueries{DB: db},
		AuthQueries: &queries.AuthQueries{DB: db},
	}

	return dbInstance, nil
}

func setupConnectionPool(db *sqlx.DB, config *viper.Viper) {
	maxOpenConns := config.GetInt("DB_MAX_CONNECTIONS")
	maxIdleConns := config.GetInt("DB_MAX_IDLE_CONNECTIONS")
	maxLifetime := config.GetInt("DB_MAX_LIFETIME_CONNECTIONS")

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
}

func CloseDBConnection() error {
	mu.Lock()
	defer mu.Unlock()

	if dbInstance != nil && dbInstance.UserQueries != nil {
		if err := dbInstance.UserQueries.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		dbInstance = nil
		fmt.Printf("Database connection closed successfully")
	}
	return nil
}

func HealthCheck() error {
	mu.RLock()
	defer mu.RUnlock()

	if dbInstance == nil || dbInstance.UserQueries == nil {
		return fmt.Errorf("database connection not initialized")
	}

	if err := dbInstance.UserQueries.DB.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}
