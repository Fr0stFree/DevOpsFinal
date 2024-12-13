package app

import (
	"net/http"

	"project_sem/internal/db"
	"project_sem/internal/handlers"
)

func NewServerRouter(repo *db.Repository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v0/prices", handlers.GetPrices(repo))
	mux.HandleFunc("POST /api/v0/prices", handlers.CreatePrices(repo))
	return mux
}
