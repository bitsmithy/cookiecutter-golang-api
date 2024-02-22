package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sync"

	"{{ cookiecutter.module_path }}/internal/env"
	"{{ cookiecutter.module_path }}/internal/log"
	"{{ cookiecutter.module_path }}/internal/version"
)

func main() {
	err := run()
	if err != nil {
		trace := string(debug.Stack())
		log.New().Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	httpPort int
}

type application struct {
	config config
	wg     sync.WaitGroup
}

func run() error {
	var cfg config

	cfg.httpPort = env.GetInt("{{ cookiecutter.module_name.upper().replace('-', '_') }}_PORT", {{ cookiecutter.server_port }})

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	app := &application{
		config: cfg,
	}

	return app.serveHTTP()
}
