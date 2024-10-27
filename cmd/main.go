package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"todolist/internal/config"
	add_task2 "todolist/internal/http-server/handlers/add-task"
	complete_task "todolist/internal/http-server/handlers/complete-task"
	"todolist/internal/http-server/middleware/auth"
	"todolist/internal/http-server/middleware/logger"
	"todolist/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.Username,
		cfg.DB.DBPassword,
		cfg.DB.Host, cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	storage, err := postgresql.New(dsn)
	if err != nil {
		log.Error("failed to init storage:", err)
	}
	defer storage.DB.Close()

	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.URLFormat,
		logger.New(log),
		auth.New(log),
	)

	router.Post("/add_task", add_task2.New(log, storage))
	router.Patch("/complete", complete_task.New(log, storage))

	serverAddr := cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port
	log.Info("starting server...", slog.String("address", serverAddr))
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
