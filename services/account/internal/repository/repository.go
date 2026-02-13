package repository

import (
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/Molov30/go-microservices/services/account/internal/config"
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

	return &Repository{
		db:     db,
		logger: logger,
	}, nil
}
