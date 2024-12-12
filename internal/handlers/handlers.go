package handlers

import (
	"net/http"
	"log"
	"encoding/json"

	"project_sem/internal/db"
	"project_sem/internal/utils"
)

func GetPrices(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prices, err := repo.GetPrices()
		if err != nil {
			log.Printf("failed to load prices: %v\n", err)
			http.Error(w, "failed to get prices", http.StatusInternalServerError)
			return
		}
		err = utils.ArchivePrices(prices, w, "data.csv")
		if err != nil {
			log.Printf("failed to save prices to csv: %v\n", err)
			http.Error(w, "failed to get prices", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
	}
}

func CreatePrices(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prices, err := utils.UnarchivePrices(r.Body)
		if err != nil {
			log.Printf("failed to load prices from inctoming file: %v\n", err)
			http.Error(w, "failed to upload prices", http.StatusInternalServerError)
			return
		}
		err = repo.CreatePrices(prices)
		if err != nil {
			log.Printf("failed to save prices into db: %v\n", err)
			http.Error(w, "failed to upload prices", http.StatusInternalServerError)
			return
		}
		stats := utils.CalculatePriceStats(prices)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(stats)
	}
}