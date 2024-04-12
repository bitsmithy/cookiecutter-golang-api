package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{ cookiecutter.module_path }}/internal/log"
	"{{ cookiecutter.module_path }}/internal/telemetry"
)

const (
	serviceName           = "{{ cookiecutter.module_name }}"
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (app *application) serveHTTP() error {
	l := log.New()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		Handler:      app.routes(),
		ErrorLog:     l.StdLogger(log.ErrorLevel),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	shutdownServerChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownServerChan <- srv.Shutdown(ctx)
	}()

	// Set up OpenTelemetry.
	otelCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := telemetry.Setup(otelCtx, serviceName)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	l.With("server.addr", srv.Addr).Info("starting server")

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownServerChan
	if err != nil {
		return err
	}

	l.With("server.addr", srv.Addr).Info("stopped server")

	app.wg.Wait()
	return nil
}
