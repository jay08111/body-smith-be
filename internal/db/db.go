package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	"body-smith-be/internal/config"
	migrationfiles "body-smith-be/migrations"
)

func New(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", cfg.DSN(true))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sqlx.DB, databaseName string) error {
	stdDB := db.DB

	driver, err := migratemysql.WithInstance(stdDB, &migratemysql.Config{
		DatabaseName: databaseName,
	})
	if err != nil {
		return fmt.Errorf("create migrate mysql driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationfiles.Files, ".")
	if err != nil {
		return fmt.Errorf("create migrate source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, databaseName, driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func EnsureDatabase(cfg *config.Config) error {
	rootDB, err := NewRootConnection(cfg)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rootDB.PingContext(ctx); err != nil {
		return err
	}

	_, err = rootDB.ExecContext(ctx, fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		cfg.DBName,
	))
	return err
}

func NewRootConnection(cfg *config.Config) (*sql.DB, error) {
	return sql.Open("mysql", cfg.DSN(false))
}
