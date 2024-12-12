package utils

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"project_sem/internal/db"
	"strconv"
)

func ArchivePrices(prices []db.Price, w io.Writer, fileName string) error {
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

func UnarchivePrices(r io.ReadCloser) ([]db.Price, error) {
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

type PriceStats struct {
	TotalItems      int `json:"total_items"`
	TotalCategories int `json:"total_categories"`
	TotalPrice      int `json:"total_price"`
}

func CalculatePriceStats(prices []db.Price) PriceStats {
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
