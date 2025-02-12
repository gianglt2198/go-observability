package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"coffee-shop-api/config"
	"coffee-shop-api/monitoring"
)

type Database interface {
	DB() *gorm.DB
	Close() error
	Ping() error
}

type database struct {
	db *gorm.DB
}

func NewDatabase(cfg *config.Config, log *monitoring.AppLogger) (Database, error) {

	logger := NewLogger(log.GetLogger()).LogMode(gormlogger.Info)

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetConnMaxLifetime(cfg.Database.Timeout)

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &database{db: db}, nil
}

func MustNewDatabase(cfg *config.Config, logger *monitoring.AppLogger) Database {
	db, err := NewDatabase(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return db
}

func (d *database) DB() *gorm.DB {
	return d.db
}

func (d *database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *database) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
