package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project_sem/internal/config"
	"project_sem/internal/db"
)

type App struct {
	server *http.Server
}

func New(config config.Config) *App {
	repo, err := db.NewRepository(config.DB)
	if err != nil {
		log.Fatalf("failed to create repository with error %s", err)
	}
	router := NewServerRouter(repo)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}
	return &App{server}
}

func (app *App) Run() {
	go func() {
		log.Printf("starting server on %s...\n", app.server.Addr)
		err := app.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server has failed with %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server shutdown failed with %s", err)
	}
	log.Println("server has been shutdown")
}
