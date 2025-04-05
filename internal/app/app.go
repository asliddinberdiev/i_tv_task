package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/asliddinberdiev/i_tv_task/internal/config"
	deliveryHttp "github.com/asliddinberdiev/i_tv_task/internal/delivery/http"
	v1 "github.com/asliddinberdiev/i_tv_task/internal/delivery/http/v1"
	"github.com/asliddinberdiev/i_tv_task/internal/modules/movie"
	"github.com/asliddinberdiev/i_tv_task/internal/modules/user"
	"github.com/asliddinberdiev/i_tv_task/internal/storage/postgres"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"go.uber.org/fx"
)

type App struct {
	fxApp *fx.App
	log   logger.Logger
}

func NewCreateApp() *App {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	log := logger.NewLogger(cfg.App.LogLevel, cfg.App.ServiceName)
	defer logger.Cleanup(log)

	app := fx.New(
		fx.Provide(func() *config.Config {
			return cfg
		}),

		fx.Provide(func() logger.Logger {
			return log
		}),

		postgres.Module,
		user.Module,
		movie.Module,
		deliveryHttp.Module,
		v1.Module,

		fx.Invoke(func(lc fx.Lifecycle, handler *deliveryHttp.Handler, cfg *config.Config, log logger.Logger, psql postgres.PostgresDB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := psql.AutoMigrate(
						&user.User{},
						&movie.Movie{},
					); err != nil {
						log.Error("failed to auto migrate", logger.Error(err))
						return err
					}

					server := &http.Server{
						Addr:         cfg.GetAppAddr(),
						Handler:      handler.Router,
						ReadTimeout:  cfg.App.ReadTimeout,
						WriteTimeout: cfg.App.WriteTimeout,
						IdleTimeout:  cfg.App.IdleTimeout,
					}

					go func() {
						if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
							log.Fatal("failed to start server", logger.Error(err))
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
					defer cancel()

					<-ctx.Done()
					return ctx.Err()
				},
			})
		}),
	)

	return &App{
		fxApp: app,
		log:   log,
	}
}

func (a *App) Start() error {
	return a.fxApp.Start(context.Background())
}

func (a *App) Stop() error {
	return a.fxApp.Stop(context.Background())
}
