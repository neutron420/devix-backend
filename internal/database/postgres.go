package database

import (
	"fmt"
	"time"

	"devix-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		PrepareStmt: true,
	}

	db, err := gorm.Open(postgres.Open(cfg.URL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(int(cfg.MinConns))
	sqlDB.SetMaxOpenConns(int(cfg.MaxConns))
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
