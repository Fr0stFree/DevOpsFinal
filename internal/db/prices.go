package db

import (
	"fmt"
	"strings"
	"time"
)

type Price struct {
	ID         int
	Name       string
	Category   string
	Price      float64
	CreateDate time.Time
}

func (r *Repository) GetPrices() ([]Price, error) {
	prices := make([]Price, 0)

	rows, err := r.db.Query("SELECT * FROM prices")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var price Price
		err = rows.Scan(&price.ID, &price.Name, &price.Category, &price.Price, &price.CreateDate)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	return prices, nil
}

func (r *Repository) CreatePrices(prices []Price) error {
	values := make([]string, len(prices))
	for idx, price := range prices {
		values[idx] = fmt.Sprintf("('%d', '%s', '%s', '%f', '%s')", price.ID, price.Name, price.Category, price.Price, price.CreateDate.Format("2006-01-02"))
	}
	_, err := r.db.Exec(fmt.Sprintf("INSERT INTO prices (id, name, category, price, create_date) VALUES %s", strings.Join(values, ",")))
	if err != nil {
		return err
	}
	return nil
}
