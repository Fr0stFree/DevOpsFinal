package handlers

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"

	"project_sem/internal/db"
)

func getPrices(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prices, err := repo.GetPrices()
		if err != nil {
			log.Printf("failed to load prices: %v\n", err)
			http.Error(w, "failed to get prices", http.StatusInternalServerError)
			return
		}
		err = archivePrices(prices, w, "data.csv")
		if err != nil {
			log.Printf("failed to save prices to csv: %v\n", err)
			http.Error(w, "failed to get prices", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
	}
}

func archivePrices(prices []db.Price, w io.Writer, fileName string) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	
	file, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	csvWriter.Write([]string{"id", "name", "category", "price", "create_date"})
	for _, price := range prices {
		record := []string{
			fmt.Sprintf("%d", price.ID),
			price.Name,
			price.Category,
			fmt.Sprintf("%.2f", price.Price),
			price.CreateDate,
		}
		err := csvWriter.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}
