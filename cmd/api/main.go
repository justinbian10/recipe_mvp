package main

import (
	"flag"
	"net/http"
	"os"
	"log/slog"
	"time"
)

type config struct {
	port int
	env string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr: ":8080",
		Handler: app.routes(),
		IdleTimeout: 5 * time.Second,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
