package resthttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"ghorkov32/proletariat-budget-be/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const serverShutdownTimeoutDuration = 60 * time.Second

type App struct {
	cfg    *config.App
	Server *http.Server
}

func NewHTTPServer(
	cfg *config.App,
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) *App {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

	r.Mount("/metrics", promhttp.Handler())
	r.With(middlewares...).Mount("/", handler)

	server := http.Server{
		Handler:     r,
		Addr:        fmt.Sprintf(":%d", cfg.ServerPort),
		ReadTimeout: cfg.ReadTimeout,
	}

	return &App{
		cfg:    cfg,
		Server: &server,
	}
}

func (s *App) Start() {
	log.Info().Int("port", s.cfg.ServerPort).Msg("starting http sever")
	err := s.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("http server closed")
	} else if err != nil {
		log.Fatal().Err(err).Msg("failed to start http server")
	}
}

func (s *App) Shutdown(ctx context.Context) error {
	log.Info().Msg("shutting down http server")
	if s.Server != nil {
		shutdownCtx, shutdownRelease := context.WithTimeout(ctx, serverShutdownTimeoutDuration)
		err := s.Server.Shutdown(shutdownCtx)
		shutdownRelease()

		return err
	}

	return nil
}
