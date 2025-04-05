package postgres

import (
	"context"

	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Module = fx.Module("postgres", fx.Provide(NewPostgres))

type PostgresDB interface {
	DB() *gorm.DB
	Close() error
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	AutoMigrate(models ...interface{}) error
}

type postgresDB struct {
	db *gorm.DB
}

func NewPostgres(cfg *config.Config) (PostgresDB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	if cfg.App.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetPostgresDSN()), gormConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql.DB from gorm")
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		return nil, errors.Wrap(err, "failed to ping postgres")
	}

	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime)

	return &postgresDB{db: db}, nil
}

func (p *postgresDB) DB() *gorm.DB {
	return p.db
}

func (p *postgresDB) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get sql.DB from gorm")
	}
	return sqlDB.Close()
}

func (p *postgresDB) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (p *postgresDB) AutoMigrate(models ...interface{}) error {
	return p.db.AutoMigrate(models...)
}
