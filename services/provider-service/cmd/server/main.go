package main

import (
	"log"

	"github.com/freedom-music/provider-service/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}
	if err := a.Run(); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
