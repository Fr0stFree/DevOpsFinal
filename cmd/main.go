package main

import (
	"project_sem/internal/app"
	"project_sem/internal/config"
)

func main() {
	cfg := config.Load()
	instance := app.New(cfg)
	instance.Run()
}
