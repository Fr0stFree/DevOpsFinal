package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project_sem/internal/config"
	"project_sem/internal/db"
	"project_sem/internal/handlers"
)

type Application struct {
	server *http.Server
}

func NewApp(config config.Config) *Application {
	repo, err := db.NewRepository(config.DB)
	if err != nil {
		log.Fatalf("failed to create repository with error %s", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v0/prices", handlers.GetPrices(repo))
	mux.HandleFunc("POST /api/v0/prices", handlers.CreatePrices(repo))
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Server.Port),
		Handler:      mux,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}
	return &Application{server}
}

func (app *Application) Run() {
	go func() {
		log.Printf("starting server on %s...\n", app.server.Addr)
		err := app.server.ListenAndServe()
		if err != nil {
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
