package main

import (
	"log"

	"github.com/freedom-music/api-gateway/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}
	if err := application.Run(); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
