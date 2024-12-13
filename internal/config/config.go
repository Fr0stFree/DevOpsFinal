package config

import (
	"log"
	"time"
)

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

func Load() Config {
	log.Println("loading configuration...")
	return Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		DB: DBConfig{
			Host:     "localhost",
			Port:     5432,
			Password: "val1dat0r",
			User:     "validator",
			Name:     "project-sem-1",
		},
	}
}
