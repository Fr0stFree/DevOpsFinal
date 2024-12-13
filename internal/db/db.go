package db

import (
	"database/sql"
	"fmt"
	"log"
	"project_sem/internal/config"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg config.DBConfig) (*Repository, error) {
	log.Println("connecting to database...")
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &Repository{}, err
	}
	err = db.Ping()
	if err != nil {
		return &Repository{}, err
	}
	log.Printf("successfully connected to database '%s'\n", cfg.Name)
	return &Repository{db}, nil
}
