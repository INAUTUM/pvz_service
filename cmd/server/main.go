package main

import (
	"pvz_service/internal/app"
)

func main() {
    cfg := app.LoadConfig()
    application := app.NewApplication(cfg) // Исправлено на NewApplication
    if err := application.Run(); err != nil {
        panic(err)
    }
}