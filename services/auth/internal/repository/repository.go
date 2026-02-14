package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Molov30/go-microservices/services/auth/internal/config"
	_ "github.com/Molov30/go-microservices/services/auth/internal/repository/migrations" // Register goose migrations via init()
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewRepository(cfg *config.Config, logger *zerolog.Logger) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(cfg.DbDSN), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	logger.Info().Msg("successfully connected to the database")

	migrationsConn, err := db.DB()
	if migrationsConn == nil {
		return nil, fmt.Errorf("db connection is nil: %w", err)
	}

	logger.Info().Msg("run migrations")
	err = runMigrations(migrationsConn)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info().Msg("migrations completed")

	return &Repository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	done := make(chan error)
	go func() {
		defer close(done)

		db, err := r.db.DB()
		if err != nil {
			done <- fmt.Errorf("failed to get db: %w", err)
			return
		}

		err = db.Close()
		if err != nil {
			done <- fmt.Errorf("failed to close db: %w", err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		return errors.New("repository forcing stop")
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to close repository: %w", err)
		}
	}
	return nil
}

func runMigrations(db *sql.DB) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to up migrations: %w", err)
	}

	return nil
}
