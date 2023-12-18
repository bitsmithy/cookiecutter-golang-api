package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rs/cors"
	"{{ cookiecutter.module_path }}/internal/util"
)

type srvConfig struct {
	port int
	env  string
}

type server struct {
	config srvConfig
	logger *slog.Logger
	wg     sync.WaitGroup //lint:ignore U1000 useful templating code, remove lint:ignore when using for the first time
}

func main() {
	handler := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	cfg := configFromEnv(logger)

	srv := &server{
		config: cfg,
		logger: logger,
	}

	err := srv.start()
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}
}

func (s *server) start() error {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposedHeaders: []string{"Content-Length"},
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.port),
		Handler:      c.Handler(s.routes()),
		ErrorLog:     slog.NewLogLogger(s.logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {
		// Do NOT use server.background() as that will add to the wait group and stall
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Blocks until signal is received
		sig := <-quit

		s.logger.Info("server shutdown requested", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// Wait for background tasks, then shutdown
		s.logger.Info("waiting for background tasks to complete")
		shutdownError <- nil
	}()

	s.logger.Info("starting server", "address", srv.Addr, "env", s.config.env)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		// http.ErrServerClosed is returned when Shutdown is called on srv.
		// This is intentional so we ignore it
		return err
	}

	err = <-shutdownError
	if err != nil {
		// If there's an issue with shutting down the server, surface the error.
		return err
	}

	s.logger.Info("stopped server")
	return nil
}

func configFromEnv(logger *slog.Logger) srvConfig {
	cfg := srvConfig{}

	port, err := strconv.Atoi(util.GetEnv("{{ cookiecutter.module_name.upper() }}_PORT", "{{ cookiecutter.server_port }}"))
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}

	cfg.port = port
	cfg.env = util.GetEnv("{{ cookiecutter.module_name.upper() }}_ENV", "development")

	return cfg
}
