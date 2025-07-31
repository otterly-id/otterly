package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/otterly-id/otterly/backend/app/queries"
	"github.com/otterly-id/otterly/backend/pkg/utils"
	"go.uber.org/zap"
)

var (
	dbInstance *Queries
	mu         sync.RWMutex
	logger     = utils.NewLogger()
)

type Queries struct {
	*queries.UserQueries
}

func GetDBConnection() (*Queries, error) {
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

	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	setupConnectionPool(db)

	dbInstance = &Queries{
		UserQueries: &queries.UserQueries{DB: db},
	}

	logger.Info("Database connection pool initialized successfully")
	return dbInstance, nil
}

func setupConnectionPool(db *sqlx.DB) {
	maxOpenConns := utils.GetEnvAsInt("DB_MAX_CONNECTIONS", 25)
	maxIdleConns := utils.GetEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5)
	maxLifetime := utils.GetEnvAsInt("DB_MAX_LIFETIME_CONNECTIONS", 300)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	logger.Info("Database connection pool configured",
		zap.Int("max_open_connections", maxOpenConns),
		zap.Int("max_idle_connections", maxIdleConns),
		zap.Int("max_lifetime_seconds", maxLifetime))
}

func CloseDBConnection() error {
	mu.Lock()
	defer mu.Unlock()

	if dbInstance != nil && dbInstance.UserQueries != nil {
		if err := dbInstance.UserQueries.DB.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.String("error", err.Error()))
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		dbInstance = nil
		logger.Info("Database connection closed successfully")
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
