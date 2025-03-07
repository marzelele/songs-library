package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	_ "songs-library/docs"
	api "songs-library/internal/api/http"
	"songs-library/internal/config"
	"songs-library/internal/respository"
	"songs-library/internal/router"
	"songs-library/internal/service"
	"songs-library/pkg/logger/sl"
	"syscall"
	"time"
)

// @title Songs Library
// @version 0.1
// @host      localhost:8080
// @description Онлайн библиотека песен
// @BasePath /api/v1
func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	db, err := respository.NewRepository(cfg.PgDsn)
	if err != nil {
		slog.Error("database init error", sl.Err(err))
		os.Exit(1)
	}

	s := service.NewService(log, db, cfg.SongsInfoAPIURL)
	h := api.NewHandler(log, s)
	r := router.NewRouter(log, h)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r.Init(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil {
			log.Error("server listen error", sl.Err(err))
			os.Exit(1)
		}
	}()

	log.Info("server started", slog.String("port", cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Info("shutting down server...")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	err = db.Close()

	if err != nil {
		log.Error("close db client error", sl.Err(err))
	}
	err = server.Shutdown(ctx)
	if err != nil {
		log.Error("shutdown server error", sl.Err(err))
		os.Exit(1)
	}
}
