package server

import (
	"context"
	"net/http"

	"github.com/asliddinberdiev/i_tv_task/internal/config"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var Module = fx.Module("server", fx.Provide(NewServer))

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	log        logger.Logger
}

func NewServer(lc fx.Lifecycle, cfg *config.Config, handler http.Handler, log logger.Logger) *Server {
	s := &Server{
		httpServer: &http.Server{
			Addr:         cfg.GetAppAddr(),
			Handler:      handler,
			ReadTimeout:  cfg.App.ReadTimeout,
			WriteTimeout: cfg.App.WriteTimeout,
			IdleTimeout:  cfg.App.IdleTimeout,
		},
		cfg: cfg,
		log: log,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			errChan := make(chan error, 1)
			go func() {
				s.log.Info("starting server", logger.Any("addr", s.httpServer.Addr))
				if err := s.run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					errChan <- errors.Wrap(err, "failed to start server")
				} else {
					errChan <- nil
				}
			}()

			select {
			case err := <-errChan:
				if err != nil {
					s.log.Error("failed to start server", logger.Error(err))
					return err
				}
			case <-ctx.Done():
				s.log.Error("context done", logger.Error(ctx.Err()))
				return ctx.Err()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Info("shutting down server gracefully")
			shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.App.GracePeriod)
			defer cancel()

			if err := s.stop(shutdownCtx); err != nil {
				s.log.Error("failed to stop server gracefully", logger.Error(err))
				return errors.Wrap(err, "failed to stop server gracefully")
			}
			s.log.Info("server stopped gracefully")
			return nil
		},
	})

	return s
}

func (s *Server) run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
