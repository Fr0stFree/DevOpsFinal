package main

import (
	"project_sem/internal/app"
	"project_sem/internal/config"
)

func main() {
	configuration := config.Load()
	application := app.NewApp(configuration)
	application.Run()
}
