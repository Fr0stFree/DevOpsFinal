package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"project_sem/internal/db"
)

type PriceStats struct {
	TotalItems      int `json:"total_items"`
	TotalCategories int `json:"total_categories"`
	TotalPrice      int `json:"total_price"`
}

func createPrices(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prices, err := unarchivePrices(r.Body)
		if err != nil {
			log.Printf("failed to load prices from incoming file: %v\n", err)
			http.Error(w, "failed to upload prices", http.StatusInternalServerError)
			return
		}
		err = repo.CreatePrices(prices)
		if err != nil {
			log.Printf("failed to save prices into db: %v\n", err)
			http.Error(w, "failed to upload prices", http.StatusInternalServerError)
			return
		}
		stats := calculatePriceStats(prices)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(stats)
	}
}

func unarchivePrices(r io.Reader) ([]db.Price, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	prices := make([]db.Price, 0)
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		csvReader := csv.NewReader(rc)
		records, err := csvReader.ReadAll()
		if err != nil {
			return nil, err
		}
		for idx, record := range records {
			if idx == 0 {
				continue
			}
			price, err := parsePrice(record)
			if err != nil {
				return nil, err
			}
			prices = append(prices, price)
		}
		rc.Close()
	}
	return prices, nil
}


func parsePrice(record []string) (db.Price, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return db.Price{}, fmt.Errorf("failed to convert id to int %v", err)
	}
	cost, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return db.Price{}, fmt.Errorf("failed to convert cost to float %v", err)
	}
	price := db.Price{
		ID:         id,
		Name:       record[1],
		Category:   record[2],
		Price:      cost,
		CreateDate: record[4],
	}
	return price, nil
}

func calculatePriceStats(prices []db.Price) PriceStats {
	stats := PriceStats{}
	categories := make(map[string]bool)
	for _, price := range prices {
		stats.TotalItems++
		categories[price.Category] = true
		stats.TotalPrice += int(price.Price)
	}
	stats.TotalCategories = len(categories)
	return stats
}
