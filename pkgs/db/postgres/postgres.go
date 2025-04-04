package db

import (
	"context"
	"time"

	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPg(cfg *config.Config, ctx context.Context) (*gorm.DB, error) {
	dsn := cfg.GetPostgresDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql.DB from gorm")
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping postgres")
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	return db, nil
}
