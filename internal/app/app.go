package app

import (
	"context"
	"errors"
	externalLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/asliddinberdiev/i_tv_task/internal/api/http"
	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/asliddinberdiev/i_tv_task/internal/modules"
	"github.com/asliddinberdiev/i_tv_task/internal/server"
	db "github.com/asliddinberdiev/i_tv_task/pkgs/db/postgres"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		externalLog.Fatalf("failed to initialize configuration: %+v\n", err)
	}

	log := logger.NewLogger(cfg.App.LogLevel, cfg.App.ServiceName)
	defer logger.Cleanup(log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	psqlConn, err := db.ConnectPg(cfg, ctx)
	if err != nil {
		log.Fatal("failed to connect to postgres", logger.Error(err))
		return
	}
	log.Info("postgres", logger.Any("DSN", cfg.GetPostgresDSN()))

	modules := modules.NewModules(cfg, log, psqlConn)

	handlers := api.NewHandler(log, cfg, modules)

	srv := server.NewServer(cfg, handlers.Init())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("failed to run server", logger.Error(err))
		}
	}()

	log.Info("server", logger.Any("ADDR", cfg.GetAppAddr()))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatal("failed to stop server", logger.Error(err))
	}
}
