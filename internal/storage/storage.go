package storage

import (
	"github.com/asliddinberdiev/i_tv_task/internal/storage/postgres"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"storage",
	fx.Provide(NewStorage),
	postgres.Module,
)

type Storage interface {
	Postgres() postgres.PostgresDB
	Close() error
}

type storage struct {
	postgres postgres.PostgresDB
}

func NewStorage(postgres postgres.PostgresDB) Storage {
	return &storage{postgres: postgres}
}

func (s *storage) Postgres() postgres.PostgresDB {
	return s.postgres
}

func (s *storage) Close() error {
	if err := s.postgres.Close(); err != nil {
		return errors.Wrap(err, "failed to close postgres")
	}
	return nil
}
