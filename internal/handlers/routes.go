package handlers

import (
	"net/http"

	"project_sem/internal/db"
)

func NewServerRouter(repo *db.Repository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v0/prices", getPrices(repo))
	mux.HandleFunc("POST /api/v0/prices", createPrices(repo))
	return mux
}
